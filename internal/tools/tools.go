// Package tools exposes a small set of sandboxed operations that generated
// code is allowed to call. Every function is rooted at a single workDir so
// that generated code can never read or write outside of it, and shell
// commands are restricted to an allowlist of binaries.
package tools

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// workDir is the sandbox root. All paths passed to the functions below are
// resolved relative to it, and resolution fails if it escapes the root.
var workDir = "."

// SetWorkDir sets the sandbox root used by every tool call. It must be
// called once during agent startup, before any generated code runs.
func SetWorkDir(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	workDir = abs
	return nil
}

// resolve turns a user-supplied relative path into an absolute path and
// verifies it stays inside workDir.
func resolve(path string) (string, error) {
	full := filepath.Join(workDir, path)
	full = filepath.Clean(full)
	rel, err := filepath.Rel(workDir, full)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes the sandbox root %q", path, workDir)
	}
	return full, nil
}

// ReadFile returns the contents of a text file relative to the sandbox root.
func ReadFile(path string) (string, error) {
	full, err := resolve(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file relative to the sandbox root, creating
// parent directories as needed.
func WriteFile(path string, content string) error {
	full, err := resolve(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0o644)
}

// ListDir returns the entries (files and directories) directly inside path.
func ListDir(path string) ([]string, error) {
	full, err := resolve(path)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(full)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}

// FileTree renders a textual tree of path up to maxDepth levels deep.
// Hidden entries (dotfiles) and common build/VCS directories are skipped.
func FileTree(path string, maxDepth int) (string, error) {
	full, err := resolve(path)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	var walk func(dir string, prefix string, depth int) error
	walk = func(dir string, prefix string, depth int) error {
		if depth > maxDepth {
			return nil
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		filtered := entries[:0]
		for _, e := range entries {
			n := e.Name()
			if strings.HasPrefix(n, ".") || n == "node_modules" || n == "vendor" {
				continue
			}
			filtered = append(filtered, e)
		}
		for i, e := range filtered {
			last := i == len(filtered)-1
			branch := "├── "
			nextPrefix := prefix + "│   "
			if last {
				branch = "└── "
				nextPrefix = prefix + "    "
			}
			name := e.Name()
			if e.IsDir() {
				name += "/"
			}
			b.WriteString(prefix + branch + name + "\n")
			if e.IsDir() {
				if err := walk(filepath.Join(dir, e.Name()), nextPrefix, depth+1); err != nil {
					return err
				}
			}
		}
		return nil
	}
	b.WriteString(filepath.Base(full) + "/\n")
	if err := walk(full, "", 1); err != nil {
		return "", err
	}
	return b.String(), nil
}

// GrepMatch is a single grep hit.
type GrepMatch struct {
	File string
	Line int
	Text string
}

// Grep searches for a regular expression pattern in every text file under
// root and returns the matching lines.
func Grep(root string, pattern string) ([]GrepMatch, error) {
	full, err := resolve(root)
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	var matches []GrepMatch
	err = filepath.WalkDir(full, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(name, ".") {
			return nil
		}
		f, err := os.Open(p)
		if err != nil {
			return nil // skip unreadable files
		}
		defer f.Close()
		rel, _ := filepath.Rel(workDir, p)
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			if re.MatchString(line) {
				matches = append(matches, GrepMatch{File: rel, Line: lineNo, Text: line})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// CountLinesByExt walks root and returns a line count per file extension.
func CountLinesByExt(root string) (map[string]int, error) {
	full, err := resolve(root)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	err = filepath.WalkDir(full, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		if len(data) == 0 {
			return nil
		}
		ext := filepath.Ext(name)
		if ext == "" {
			ext = "(no ext)"
		}
		n := bytes.Count(data, []byte("\n"))
		if data[len(data)-1] != '\n' {
			n++ // count a final line with no trailing newline
		}
		counts[ext] += n
		return nil
	})
	if err != nil {
		return nil, err
	}
	return counts, nil
}

// allowedCommands is the fixed allowlist of binaries generated code may
// invoke via RunCommand. Anything else is rejected before exec.
var allowedCommands = map[string]bool{
	"go": true, "git": true, "ls": true, "wc": true,
	"find": true, "cat": true, "echo": true, "gofmt": true,
}

// RunCommand runs an allowlisted command with the given arguments inside the
// sandbox root and returns combined stdout+stderr. It is never passed
// through a shell, so shell metacharacters in args have no special meaning.
func RunCommand(name string, args ...string) (string, error) {
	if !allowedCommands[name] {
		return "", fmt.Errorf("command %q is not allowed (allowed: go, git, ls, wc, find, cat, echo, gofmt)", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = workDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
