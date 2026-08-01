package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := SetWorkDir(dir); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	return dir
}

func TestReadWriteFile(t *testing.T) {
	setup(t)
	if err := WriteFile("sub/dir/hello.txt", "hi"); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := ReadFile("sub/dir/hello.txt")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got != "hi" {
		t.Fatalf("got %q, want %q", got, "hi")
	}
}

func TestResolveBlocksEscape(t *testing.T) {
	setup(t)
	if _, err := ReadFile("../outside.txt"); err == nil {
		t.Fatalf("expected error escaping sandbox")
	}
	if _, err := ReadFile("/etc/passwd"); err == nil {
		t.Fatalf("expected error for absolute path outside sandbox")
	}
}

func TestListDir(t *testing.T) {
	dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := ListDir(".")
	if err != nil {
		t.Fatalf("ListDir: %v", err)
	}
	want := map[string]bool{"a.txt": true, "sub/": true}
	if len(entries) != 2 {
		t.Fatalf("got %v", entries)
	}
	for _, e := range entries {
		if !want[e] {
			t.Fatalf("unexpected entry %q in %v", e, entries)
		}
	}
}

func TestGrep(t *testing.T) {
	dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "code.go"), []byte("package main\n// TODO: fix this\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	matches, err := Grep(".", "TODO")
	if err != nil {
		t.Fatalf("Grep: %v", err)
	}
	if len(matches) != 1 || matches[0].Line != 2 {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

func TestCountLinesByExt(t *testing.T) {
	dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("l1\nl2\nl3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.go"), []byte("l1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	counts, err := CountLinesByExt(".")
	if err != nil {
		t.Fatalf("CountLinesByExt: %v", err)
	}
	if counts[".go"] != 4 {
		t.Fatalf("got %v", counts)
	}
}

func TestCountLinesByExtSkipsBuildOutputDirs(t *testing.T) {
	dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("l1\nl2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Simulate a compiled binary artifact: arbitrary bytes that happen to
	// contain many newline bytes, which would otherwise wildly inflate the
	// line count for extension-less files.
	fakeBinary := make([]byte, 4096)
	for i := range fakeBinary {
		fakeBinary[i] = byte(i % 251)
	}
	if err := os.WriteFile(filepath.Join(dir, "bin", "agent"), fakeBinary, 0o755); err != nil {
		t.Fatal(err)
	}
	counts, err := CountLinesByExt(".")
	if err != nil {
		t.Fatalf("CountLinesByExt: %v", err)
	}
	if counts[".go"] != 2 {
		t.Fatalf("got %v", counts)
	}
	if _, ok := counts["(no ext)"]; ok {
		t.Fatalf("expected bin/ to be skipped entirely, got counts: %v", counts)
	}
}

func TestRunCommandAllowlist(t *testing.T) {
	setup(t)
	if _, err := RunCommand("rm", "-rf", "/"); err == nil {
		t.Fatalf("expected disallowed command to be rejected")
	}
	out, err := RunCommand("echo", "hello")
	if err != nil {
		t.Fatalf("RunCommand echo: %v", err)
	}
	if out != "hello\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFileTree(t *testing.T) {
	dir := setup(t)
	if err := os.MkdirAll(filepath.Join(dir, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "pkg", "f.go"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree, err := FileTree(".", 2)
	if err != nil {
		t.Fatalf("FileTree: %v", err)
	}
	if tree == "" {
		t.Fatalf("expected non-empty tree")
	}
}
