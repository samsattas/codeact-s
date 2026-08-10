package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	if err := tools.SetWorkDir(t.TempDir()); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	return &server{
		provider:    &fakeProvider{responses: responses},
		exec:        executor.New(),
		maxAttempts: 3,
		workDir:     tools.WorkDir(),
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
	if len(lines) < 2 {
		t.Fatalf("expected at least a step and a done message, got %d lines: %q", len(lines), rec.Body.String())
	}

	var step streamMessage
	if err := json.Unmarshal([]byte(lines[0]), &step); err != nil {
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
