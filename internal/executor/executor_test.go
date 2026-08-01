package executor

import (
	"context"
	"strings"
	"testing"

	"codeact-agent/internal/tools"
)

func TestExecutorRunsGeneratedCode(t *testing.T) {
	dir := t.TempDir()
	if err := tools.SetWorkDir(dir); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	exec := New()

	src := `
package main

import (
	"fmt"
	"tools"
)

func Run() (string, error) {
	if err := tools.WriteFile("hello.txt", "hi there"); err != nil {
		return "", err
	}
	content, err := tools.ReadFile("hello.txt")
	if err != nil {
		return "", err
	}
	fmt.Println("wrote and read back a file")
	return content, nil
}
`
	res, err := exec.Run(context.Background(), src)
	if err != nil {
		t.Fatalf("Run returned executor error: %v", err)
	}
	if res.ToolError != "" {
		t.Fatalf("Run returned tool error: %s", res.ToolError)
	}
	if res.ReturnValue != "hi there" {
		t.Fatalf("unexpected return value: %q", res.ReturnValue)
	}
	if !strings.Contains(res.Stdout, "wrote and read back a file") {
		t.Fatalf("expected stdout to contain print output, got: %q", res.Stdout)
	}
}

func TestExecutorSurfacesToolError(t *testing.T) {
	dir := t.TempDir()
	if err := tools.SetWorkDir(dir); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	exec := New()

	src := `
package main

import "tools"

func Run() (string, error) {
	return tools.ReadFile("does-not-exist.txt")
}
`
	res, err := exec.Run(context.Background(), src)
	if err != nil {
		t.Fatalf("Run returned executor error: %v", err)
	}
	if res.ToolError == "" {
		t.Fatalf("expected a tool error for missing file, got none")
	}
}

func TestExecutorSurfacesCompileError(t *testing.T) {
	dir := t.TempDir()
	if err := tools.SetWorkDir(dir); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	exec := New()

	src := `
package main

func Run() (string, error) {
	this is not valid go
}
`
	_, err := exec.Run(context.Background(), src)
	if err == nil {
		t.Fatalf("expected a compile error, got none")
	}
}

func TestStdlibOSPackageIsUnavailable(t *testing.T) {
	dir := t.TempDir()
	if err := tools.SetWorkDir(dir); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	exec := New()

	// Generated code must not be able to bypass the tools sandbox by
	// importing os directly.
	src := `
package main

import "os"

func Run() (string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}
	return string(rune(len(entries))), nil
}
`
	_, err := exec.Run(context.Background(), src)
	if err == nil {
		t.Fatalf("expected importing os to fail to compile, it succeeded")
	}
}

func TestSandboxBlocksPathEscape(t *testing.T) {
	dir := t.TempDir()
	if err := tools.SetWorkDir(dir); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	exec := New()

	src := `
package main

import "tools"

func Run() (string, error) {
	return tools.ReadFile("../../etc/passwd")
}
`
	res, err := exec.Run(context.Background(), src)
	if err != nil {
		t.Fatalf("Run returned executor error: %v", err)
	}
	if res.ToolError == "" || !strings.Contains(res.ToolError, "escapes the sandbox") {
		t.Fatalf("expected sandbox escape error, got: %q", res.ToolError)
	}
}
