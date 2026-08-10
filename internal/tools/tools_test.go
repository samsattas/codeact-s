package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) (*Sandbox, string) {
	t.Helper()
	dir := t.TempDir()
	sb, err := NewSandbox(dir)
	if err != nil {
		t.Fatalf("NewSandbox: %v", err)
	}
	return sb, dir
}

func TestReadWriteFile(t *testing.T) {
	sb, _ := setup(t)
	if err := sb.WriteFile("sub/dir/hello.txt", "hi"); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := sb.ReadFile("sub/dir/hello.txt")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got != "hi" {
		t.Fatalf("got %q, want %q", got, "hi")
	}
}

func TestResolveBlocksEscape(t *testing.T) {
	sb, _ := setup(t)
	if _, err := sb.ReadFile("../outside.txt"); err == nil {
		t.Fatalf("expected error escaping sandbox")
	}
	if _, err := sb.ReadFile("/etc/passwd"); err == nil {
		t.Fatalf("expected error for absolute path outside sandbox")
	}
}

func TestListDir(t *testing.T) {
	sb, dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := sb.ListDir(".")
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
	sb, dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "code.go"), []byte("package main\n// TODO: fix this\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	matches, err := sb.Grep(".", "TODO")
	if err != nil {
		t.Fatalf("Grep: %v", err)
	}
	if len(matches) != 1 || matches[0].Line != 2 {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

func TestCountLinesByExt(t *testing.T) {
	sb, dir := setup(t)
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("l1\nl2\nl3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.go"), []byte("l1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	counts, err := sb.CountLinesByExt(".")
	if err != nil {
		t.Fatalf("CountLinesByExt: %v", err)
	}
	if counts[".go"] != 4 {
		t.Fatalf("got %v", counts)
	}
}

func TestCountLinesByExtSkipsBuildOutputDirs(t *testing.T) {
	sb, dir := setup(t)
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
	counts, err := sb.CountLinesByExt(".")
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
	sb, _ := setup(t)
	if _, err := sb.RunCommand("rm", "-rf", "/"); err == nil {
		t.Fatalf("expected disallowed command to be rejected")
	}
	out, err := sb.RunCommand("echo", "hello")
	if err != nil {
		t.Fatalf("RunCommand echo: %v", err)
	}
	if out != "hello\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFileTree(t *testing.T) {
	sb, dir := setup(t)
	if err := os.MkdirAll(filepath.Join(dir, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "pkg", "f.go"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree, err := sb.FileTree(".", 2)
	if err != nil {
		t.Fatalf("FileTree: %v", err)
	}
	if tree == "" {
		t.Fatalf("expected non-empty tree")
	}
}

func TestNewSandboxRejectsNonDirectory(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "not-a-dir.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewSandbox(file); err == nil {
		t.Fatalf("expected NewSandbox to reject a non-directory path")
	}
}

func TestNewSandboxRejectsMissingDirectory(t *testing.T) {
	if _, err := NewSandbox(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Fatalf("expected NewSandbox to reject a missing directory")
	}
}

func TestSandboxesAreIndependent(t *testing.T) {
	sbA, _ := setup(t)
	sbB, _ := setup(t)
	if err := sbA.WriteFile("only-in-a.txt", "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := sbB.ReadFile("only-in-a.txt"); err == nil {
		t.Fatalf("expected sandbox B to not see files written to sandbox A's root")
	}
}
