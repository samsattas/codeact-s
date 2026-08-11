package llm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFromEnvPrefersAnthropicWhenKeySet(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	t.Setenv("ANTHROPIC_MODEL", "")
	t.Setenv("OLLAMA_HOST", "")

	p, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if got, want := p.Name(), "anthropic:claude-sonnet-4-5"; got != want {
		t.Fatalf("Name() = %q, want %q (default model)", got, want)
	}
}

func TestFromEnvRespectsAnthropicModelOverride(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	t.Setenv("ANTHROPIC_MODEL", "claude-haiku-4-5")

	p, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if got, want := p.Name(), "anthropic:claude-haiku-4-5"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestFromEnvFallsBackToOllamaWhenReachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OLLAMA_HOST", srv.URL)
	t.Setenv("OLLAMA_MODEL", "qwen2.5-coder:7b")

	p, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if got, want := p.Name(), "ollama:qwen2.5-coder:7b"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestFromEnvErrorsWhenOllamaUnreachable(t *testing.T) {
	unreachable := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	addr := unreachable.URL
	unreachable.Close() // closed immediately, so the address is now refusing connections

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OLLAMA_HOST", addr)

	_, err := FromEnv()
	if err == nil {
		t.Fatalf("expected an error when neither Anthropic nor Ollama is configured/reachable")
	}
	if !strings.Contains(err.Error(), "ollama") {
		t.Fatalf("expected error to mention ollama, got: %v", err)
	}
}
