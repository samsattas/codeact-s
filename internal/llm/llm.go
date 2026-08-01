// Package llm defines a minimal provider interface for generating code, plus
// implementations backed by a local Ollama server or the Anthropic API.
package llm

import (
	"context"
	"fmt"
	"os"
)

// Provider generates a completion for a system + user prompt pair. It is the
// only seam the agent depends on, so swapping models or vendors never
// touches the agent loop or the executor.
type Provider interface {
	// Generate returns the raw model output for the given prompts.
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	// Name identifies the provider for logging/debugging.
	Name() string
}

// FromEnv picks a provider based on environment variables:
//   - If ANTHROPIC_API_KEY is set, use Anthropic (model from ANTHROPIC_MODEL,
//     default "claude-sonnet-4-5").
//   - Otherwise, fall back to a local Ollama server (host from OLLAMA_HOST,
//     default "http://localhost:11434"; model from OLLAMA_MODEL, default
//     "qwen2.5-coder:7b").
func FromEnv() (Provider, error) {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		model := os.Getenv("ANTHROPIC_MODEL")
		if model == "" {
			model = "claude-sonnet-4-5"
		}
		return NewAnthropicProvider(key, model), nil
	}

	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "qwen2.5-coder:7b"
	}
	p := NewOllamaProvider(host, model)
	if err := p.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("no ANTHROPIC_API_KEY set and Ollama at %s is unreachable (%w); "+
			"start Ollama with `ollama serve` and pull a model with `ollama pull %s`, "+
			"or set ANTHROPIC_API_KEY to use Claude instead", host, err, model)
	}
	return p, nil
}
