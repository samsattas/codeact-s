package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestAnthropicProvider(baseURL string) *AnthropicProvider {
	p := NewAnthropicProvider("test-key", "claude-sonnet-4-5")
	p.baseURL = baseURL
	return p
}

func TestAnthropicProviderName(t *testing.T) {
	p := NewAnthropicProvider("test-key", "claude-sonnet-4-5")
	if got, want := p.Name(), "anthropic:claude-sonnet-4-5"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestAnthropicProviderGenerateSuccess(t *testing.T) {
	var gotReq anthropicRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "test-key" {
			t.Fatalf("x-api-key header = %q, want %q", got, "test-key")
		}
		if got := r.Header.Get("anthropic-version"); got == "" {
			t.Fatalf("expected an anthropic-version header")
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"content": []map[string]string{
				{"type": "text", "text": "```go\npackage main\n```"},
			},
		})
	}))
	defer srv.Close()

	p := newTestAnthropicProvider(srv.URL)
	out, err := p.Generate(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if out != "```go\npackage main\n```" {
		t.Fatalf("unexpected output: %q", out)
	}

	if gotReq.Model != "claude-sonnet-4-5" {
		t.Fatalf("request model = %q, want %q", gotReq.Model, "claude-sonnet-4-5")
	}
	if gotReq.System != "system prompt" {
		t.Fatalf("request system = %q, want %q", gotReq.System, "system prompt")
	}
	if len(gotReq.Messages) != 1 || gotReq.Messages[0].Role != "user" || gotReq.Messages[0].Content != "user prompt" {
		t.Fatalf("unexpected messages: %+v", gotReq.Messages)
	}
}

func TestAnthropicProviderGenerateConcatenatesTextBlocks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"content": []map[string]string{
				{"type": "text", "text": "part one "},
				{"type": "text", "text": "part two"},
			},
		})
	}))
	defer srv.Close()

	p := newTestAnthropicProvider(srv.URL)
	out, err := p.Generate(context.Background(), "sys", "user")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if out != "part one part two" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestAnthropicProviderGenerateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"type": "invalid_request_error", "message": "bad api key"},
		})
	}))
	defer srv.Close()

	p := newTestAnthropicProvider(srv.URL)
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil || !strings.Contains(err.Error(), "bad api key") {
		t.Fatalf("expected an error mentioning the anthropic error message, got: %v", err)
	}
}

func TestAnthropicProviderGenerateHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	p := newTestAnthropicProvider(srv.URL)
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil {
		t.Fatalf("expected an error for a non-200 response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected error to mention the status code, got: %v", err)
	}
}
