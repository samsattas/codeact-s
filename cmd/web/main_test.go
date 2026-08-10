package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codeact-agent/internal/executor"
	"codeact-agent/internal/tools"
)

// fakeProvider scripts one canned response per call, so the HTTP handler
// tests exercise the real generate/execute/fix loop deterministically,
// without a network call to Ollama or Anthropic.
type fakeProvider struct {
	responses []string
	calls     int
}

func (f *fakeProvider) Name() string { return "fake" }

func (f *fakeProvider) Generate(_ context.Context, _ string, _ string) (string, error) {
	if f.calls >= len(f.responses) {
		f.calls++
		return "", nil
	}
	r := f.responses[f.calls]
	f.calls++
	return r, nil
}

func newTestServer(t *testing.T, responses []string) *server {
	t.Helper()
	sb, err := tools.NewSandbox(t.TempDir())
	if err != nil {
		t.Fatalf("NewSandbox: %v", err)
	}
	return &server{
		provider:    &fakeProvider{responses: responses},
		exec:        executor.New(),
		maxAttempts: 3,
		defaultDir:  sb.Root(),
	}
}

func TestHandleInfo(t *testing.T) {
	srv := newTestServer(t, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	rec := httptest.NewRecorder()
	srv.handleInfo(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var info struct {
		Provider    string `json:"provider"`
		Workdir     string `json:"workdir"`
		MaxAttempts int    `json:"maxAttempts"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if info.Provider != "fake" || info.MaxAttempts != 3 || info.Workdir == "" {
		t.Fatalf("unexpected info payload: %+v", info)
	}
}

func TestHandleRunRejectsNonPost(t *testing.T) {
	srv := newTestServer(t, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/run", nil)
	rec := httptest.NewRecorder()
	srv.handleRun(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestHandleRunRejectsMissingTask(t *testing.T) {
	srv := newTestServer(t, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewReader([]byte(`{}`)))
	rec := httptest.NewRecorder()
	srv.handleRun(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandleRunRejectsInvalidDir(t *testing.T) {
	srv := newTestServer(t, nil)

	body, _ := json.Marshal(map[string]string{"task": "do something", "dir": filepath.Join(t.TempDir(), "does-not-exist")})
	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.handleRun(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body: %s", rec.Code, rec.Body.String())
	}
}

const validRunSnippet = `package main

import "tools"

func Run() (string, error) {
	if err := tools.WriteFile("note.txt", "hi"); err != nil {
		return "", err
	}
	return "wrote note.txt", nil
}
`

func TestHandleRunStreamsSuccess(t *testing.T) {
	srv := newTestServer(t, []string{"```go\n" + validRunSnippet + "```"})

	body, _ := json.Marshal(map[string]string{"task": "write a note"})
	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.handleRun(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/x-ndjson" {
		t.Fatalf("content-type = %q, want application/x-ndjson", ct)
	}

	lines := strings.Split(strings.TrimSpace(rec.Body.String()), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected a start, step, and done message, got %d lines: %q", len(lines), rec.Body.String())
	}

	var start streamMessage
	if err := json.Unmarshal([]byte(lines[0]), &start); err != nil {
		t.Fatalf("decoding start message: %v", err)
	}
	if start.Type != "start" || start.Dir != srv.defaultDir {
		t.Fatalf("unexpected start message: %+v (want dir %q)", start, srv.defaultDir)
	}

	var step streamMessage
	if err := json.Unmarshal([]byte(lines[1]), &step); err != nil {
		t.Fatalf("decoding step message: %v", err)
	}
	if step.Type != "step" || step.Attempt != 1 || step.Code == "" {
		t.Fatalf("unexpected step message: %+v", step)
	}

	var done streamMessage
	last := lines[len(lines)-1]
	if err := json.Unmarshal([]byte(last), &done); err != nil {
		t.Fatalf("decoding done message: %v", err)
	}
	if done.Type != "done" || !done.Success || done.FinalAnswer != "wrote note.txt" {
		t.Fatalf("unexpected done message: %+v", done)
	}
}

func TestHandleRunUsesRequestDirOverDefault(t *testing.T) {
	srv := newTestServer(t, []string{"```go\n" + validRunSnippet + "```"})
	customDir := t.TempDir()

	body, _ := json.Marshal(map[string]string{"task": "write a note", "dir": customDir})
	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.handleRun(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", rec.Code, rec.Body.String())
	}

	lines := strings.Split(strings.TrimSpace(rec.Body.String()), "\n")
	var start streamMessage
	if err := json.Unmarshal([]byte(lines[0]), &start); err != nil {
		t.Fatalf("decoding start message: %v", err)
	}
	if start.Dir != customDir {
		t.Fatalf("start.Dir = %q, want %q", start.Dir, customDir)
	}

	if _, err := os.Stat(filepath.Join(customDir, "note.txt")); err != nil {
		t.Fatalf("expected note.txt to be written under the request's dir, not the server default: %v", err)
	}
	if _, err := os.Stat(filepath.Join(srv.defaultDir, "note.txt")); err == nil {
		t.Fatalf("note.txt was written under the server's default dir instead of the request's dir")
	}
}

func TestHandleBrowseListsSubdirectoriesOnly(t *testing.T) {
	srv := newTestServer(t, nil)
	root := srv.defaultDir

	if err := os.Mkdir(filepath.Join(root, "sub-b"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "sub-a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".hidden"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "file.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/browse?path="+root, nil)
	rec := httptest.NewRecorder()
	srv.handleBrowse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Path    string        `json:"path"`
		Parent  string        `json:"parent"`
		Entries []browseEntry `json:"entries"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp.Path != root {
		t.Fatalf("path = %q, want %q", resp.Path, root)
	}
	if resp.Parent != filepath.Dir(root) {
		t.Fatalf("parent = %q, want %q", resp.Parent, filepath.Dir(root))
	}
	if len(resp.Entries) != 2 || resp.Entries[0].Name != "sub-a" || resp.Entries[1].Name != "sub-b" {
		t.Fatalf("unexpected entries (want only sub-a, sub-b, sorted, no dotfiles/files): %+v", resp.Entries)
	}
	if resp.Entries[0].Path != filepath.Join(root, "sub-a") {
		t.Fatalf("entry path = %q, want an absolute child path", resp.Entries[0].Path)
	}
}

func TestHandleBrowseRejectsMissingPath(t *testing.T) {
	srv := newTestServer(t, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/browse?path="+filepath.Join(t.TempDir(), "nope"), nil)
	rec := httptest.NewRecorder()
	srv.handleBrowse(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandleRunStreamsFailureAfterAttempts(t *testing.T) {
	broken := "```go\npackage main\n\nfunc Run() (string, error) {\n\tthis is not go\n}\n```"
	srv := newTestServer(t, []string{broken, broken, broken})

	body, _ := json.Marshal(map[string]string{"task": "do something impossible"})
	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.handleRun(rec, req)

	lines := strings.Split(strings.TrimSpace(rec.Body.String()), "\n")
	var done streamMessage
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &done); err != nil {
		t.Fatalf("decoding done message: %v", err)
	}
	if done.Type != "done" || done.Success {
		t.Fatalf("expected a failed done message, got: %+v", done)
	}
	if done.Error == "" {
		t.Fatalf("expected an error message on failure")
	}
}
