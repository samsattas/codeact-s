// Package tools exposes a small set of sandboxed operations that generated
// code is allowed to call. Every function is a method on a Sandbox rooted at
// a single directory so that generated code can never read or write outside
// of it, and shell commands are restricted to an allowlist of binaries. Each
// Sandbox is independent, so a caller running multiple tasks concurrently
// (e.g. the web UI, one request per directory) can give each one its own
// root without the operations interfering with each other.
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

// Sandbox confines a set of file/command operations to a single root
// directory. It is safe for concurrent use by multiple goroutines, since all
// state (just the root) is set once at construction.
type Sandbox struct {
	root string
}

// NewSandbox builds a Sandbox rooted at dir, which must exist and be a
// directory. All paths passed to the methods below are resolved relative to
// the returned root, and resolution fails if it would escape the root.
func NewSandbox(dir string) (*Sandbox, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("workdir %q: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("workdir %q is not a directory", dir)
	}
	return &Sandbox{root: abs}, nil
}

// Root returns the resolved absolute sandbox root.
func (s *Sandbox) Root() string { return s.root }

// skipDirNames lists directories that walks (Grep, CountLinesByExt,
// FileTree) skip entirely: VCS/dependency directories and common build
// output directories, whose contents are either irrelevant or, worse,
// binary artifacts that would otherwise get miscounted as text.
var skipDirNames = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	"bin":          true,
	"dist":         true,
	"build":        true,
}

func skipDir(name string) bool {
	return strings.HasPrefix(name, ".") || skipDirNames[name]
}

// resolve turns a user-supplied relative path into an absolute path and
// verifies it stays inside the sandbox root.
func (s *Sandbox) resolve(path string) (string, error) {
	full := filepath.Join(s.root, path)
	full = filepath.Clean(full)
	rel, err := filepath.Rel(s.root, full)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes the sandbox root %q", path, s.root)
	}
	return full, nil
}

// ReadFile returns the contents of a text file relative to the sandbox root.
func (s *Sandbox) ReadFile(path string) (string, error) {
	full, err := s.resolve(path)
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
func (s *Sandbox) WriteFile(path string, content string) error {
	full, err := s.resolve(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0o644)
}

// ListDir returns the entries (files and directories) directly inside path.
func (s *Sandbox) ListDir(path string) ([]string, error) {
	full, err := s.resolve(path)
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
func (s *Sandbox) FileTree(path string, maxDepth int) (string, error) {
	full, err := s.resolve(path)
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
			if skipDir(n) {
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
func (s *Sandbox) Grep(root string, pattern string) ([]GrepMatch, error) {
	full, err := s.resolve(root)
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
			if skipDir(name) {
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
		rel, _ := filepath.Rel(s.root, p)
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
func (s *Sandbox) CountLinesByExt(root string) (map[string]int, error) {
	full, err := s.resolve(root)
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
			if skipDir(name) {
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
func (s *Sandbox) RunCommand(name string, args ...string) (string, error) {
	if !allowedCommands[name] {
		return "", fmt.Errorf("command %q is not allowed (allowed: go, git, ls, wc, find, cat, echo, gofmt)", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = s.root
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
