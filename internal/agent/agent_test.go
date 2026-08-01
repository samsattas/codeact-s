package agent

import (
	"context"
	"strings"
	"testing"

	"codeact-agent/internal/executor"
	"codeact-agent/internal/tools"
)

// fakeProvider scripts a sequence of canned responses, one per call, so
// tests can exercise the retry loop deterministically without a real model.
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

func setupExecutor(t *testing.T) *executor.Executor {
	t.Helper()
	if err := tools.SetWorkDir(t.TempDir()); err != nil {
		t.Fatalf("SetWorkDir: %v", err)
	}
	return executor.New()
}

const validSnippet = "```go\n" + `package main

import "tools"

func Run() (string, error) {
	if err := tools.WriteFile("note.txt", "done"); err != nil {
		return "", err
	}
	return "wrote note.txt", nil
}
` + "\n```"

func TestAgentSucceedsFirstTry(t *testing.T) {
	exec := setupExecutor(t)
	provider := &fakeProvider{responses: []string{validSnippet}}
	ag := New(provider, exec, 3, nil)

	outcome, err := ag.Do(context.Background(), "write a note")
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if !outcome.Success {
		t.Fatalf("expected success, got failure, steps: %+v", outcome.Steps)
	}
	if outcome.FinalAnswer != "wrote note.txt" {
		t.Fatalf("unexpected final answer: %q", outcome.FinalAnswer)
	}
	if len(outcome.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(outcome.Steps))
	}
}

func TestAgentRetriesAfterCompileErrorThenSucceeds(t *testing.T) {
	exec := setupExecutor(t)
	broken := "```go\npackage main\n\nfunc Run() (string, error) {\n\tthis is not go\n}\n```"
	provider := &fakeProvider{responses: []string{broken, validSnippet}}
	ag := New(provider, exec, 3, nil)

	outcome, err := ag.Do(context.Background(), "write a note")
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if !outcome.Success {
		t.Fatalf("expected eventual success, got failure, steps: %+v", outcome.Steps)
	}
	if len(outcome.Steps) != 2 {
		t.Fatalf("expected 2 steps (fail then succeed), got %d", len(outcome.Steps))
	}
	if outcome.Steps[0].ExecErr == nil {
		t.Fatalf("expected first attempt to have an executor error")
	}
	if provider.calls != 2 {
		t.Fatalf("expected the model to be called twice, got %d", provider.calls)
	}
}

func TestAgentRetriesAfterToolErrorThenSucceeds(t *testing.T) {
	exec := setupExecutor(t)
	failing := "```go\npackage main\n\nimport \"tools\"\n\nfunc Run() (string, error) {\n\treturn tools.ReadFile(\"missing.txt\")\n}\n```"
	provider := &fakeProvider{responses: []string{failing, validSnippet}}
	ag := New(provider, exec, 3, nil)

	outcome, err := ag.Do(context.Background(), "read a file that may not exist, then write a note")
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if !outcome.Success {
		t.Fatalf("expected eventual success, got failure, steps: %+v", outcome.Steps)
	}
	if outcome.Steps[0].Output.ToolError == "" {
		t.Fatalf("expected first attempt to report a tool error")
	}
}

func TestAgentExhaustsAttempts(t *testing.T) {
	exec := setupExecutor(t)
	broken := "```go\npackage main\n\nfunc Run() (string, error) {\n\tstill not go\n}\n```"
	provider := &fakeProvider{responses: []string{broken, broken, broken}}
	ag := New(provider, exec, 3, nil)

	outcome, err := ag.Do(context.Background(), "impossible task")
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if outcome.Success {
		t.Fatalf("expected failure after exhausting attempts")
	}
	if len(outcome.Steps) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(outcome.Steps))
	}
}

func TestExtractCodeFromFence(t *testing.T) {
	raw := "Here is the code:\n```go\npackage main\n\nfunc Run() (string, error) { return \"ok\", nil }\n```\nDone."
	code, err := extractCode(raw)
	if err != nil {
		t.Fatalf("extractCode: %v", err)
	}
	if !strings.HasPrefix(code, "package main") {
		t.Fatalf("unexpected extracted code: %q", code)
	}
}

func TestExtractCodeNoFenceFails(t *testing.T) {
	_, err := extractCode("I refuse to write code.")
	if err == nil {
		t.Fatalf("expected an error when no code block is present")
	}
}
