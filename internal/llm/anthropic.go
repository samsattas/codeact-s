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

// AnthropicProvider talks to the Anthropic Messages API.
type AnthropicProvider struct {
	apiKey string
	model  string
	// baseURL defaults to the real Anthropic API; overridden by tests to
	// point at an httptest.Server instead.
	baseURL string
	client  *http.Client
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.anthropic.com",
		client:  &http.Client{Timeout: 2 * time.Minute},
	}
}

func (p *AnthropicProvider) Name() string { return fmt.Sprintf("anthropic:%s", p.model) }

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (p *AnthropicProvider) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := anthropicRequest{
		Model:     p.model,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: userPrompt},
		},
	}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling anthropic api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out anthropicResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decoding anthropic response (status %d): %w", resp.StatusCode, err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("anthropic error: %s", out.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic returned status %d: %s", resp.StatusCode, string(body))
	}

	var text string
	for _, block := range out.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return text, nil
}
