package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaProvider talks to a local Ollama server's /api/chat endpoint.
type OllamaProvider struct {
	host   string
	model  string
	client *http.Client
}

func NewOllamaProvider(host, model string) *OllamaProvider {
	return &OllamaProvider{
		host:   host,
		model:  model,
		client: &http.Client{Timeout: 2 * time.Minute},
	}
}

func (p *OllamaProvider) Name() string { return fmt.Sprintf("ollama:%s", p.model) }

// Ping checks that the Ollama server is reachable.
func (p *OllamaProvider) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.host+"/api/tags", nil)
	if err != nil {
		return err
	}
	resp, err := (&http.Client{Timeout: 3 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatResponse struct {
	Message ollamaMessage `json:"message"`
	Error   string        `json:"error"`
}

func (p *OllamaProvider) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := ollamaChatRequest{
		Model: p.model,
		Messages: []ollamaMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.host+"/api/chat", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling ollama at %s: %w", p.host, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var out ollamaChatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decoding ollama response: %w", err)
	}
	if out.Error != "" {
		return "", fmt.Errorf("ollama error: %s", out.Error)
	}
	return out.Message.Content, nil
}
