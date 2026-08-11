package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllamaProviderName(t *testing.T) {
	p := NewOllamaProvider("http://localhost:11434", "qwen2.5-coder:7b")
	if got, want := p.Name(), "ollama:qwen2.5-coder:7b"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestOllamaProviderGenerateSuccess(t *testing.T) {
	var gotReq ollamaChatRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(ollamaChatResponse{
			Message: ollamaMessage{Role: "assistant", Content: "```go\npackage main\n```"},
		})
	}))
	defer srv.Close()

	p := NewOllamaProvider(srv.URL, "qwen2.5-coder:7b")
	out, err := p.Generate(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if out != "```go\npackage main\n```" {
		t.Fatalf("unexpected output: %q", out)
	}

	if gotReq.Model != "qwen2.5-coder:7b" {
		t.Fatalf("request model = %q, want %q", gotReq.Model, "qwen2.5-coder:7b")
	}
	if gotReq.Stream {
		t.Fatalf("expected Stream: false, request must not ask for a streamed response")
	}
	if len(gotReq.Messages) != 2 || gotReq.Messages[0].Role != "system" || gotReq.Messages[1].Role != "user" {
		t.Fatalf("unexpected messages: %+v", gotReq.Messages)
	}
	if gotReq.Messages[0].Content != "system prompt" || gotReq.Messages[1].Content != "user prompt" {
		t.Fatalf("unexpected message content: %+v", gotReq.Messages)
	}
}

func TestOllamaProviderGenerateHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	p := NewOllamaProvider(srv.URL, "qwen2.5-coder:7b")
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil {
		t.Fatalf("expected an error for a non-200 response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected error to mention the status code, got: %v", err)
	}
}

func TestOllamaProviderGenerateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ollamaChatResponse{Error: "model not found"})
	}))
	defer srv.Close()

	p := NewOllamaProvider(srv.URL, "does-not-exist")
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil || !strings.Contains(err.Error(), "model not found") {
		t.Fatalf("expected an error mentioning the ollama error message, got: %v", err)
	}
}

func TestOllamaProviderPingSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	p := NewOllamaProvider(srv.URL, "qwen2.5-coder:7b")
	if err := p.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestOllamaProviderPingFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	p := NewOllamaProvider(srv.URL, "qwen2.5-coder:7b")
	if err := p.Ping(context.Background()); err == nil {
		t.Fatalf("expected an error, got none")
	}
}
