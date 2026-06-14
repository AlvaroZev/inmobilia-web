package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

type OpenAIParser struct {
	config OpenAIConfig
	client *http.Client
}

func NewOpenAIParser(config OpenAIConfig) *OpenAIParser {
	return &OpenAIParser{
		config: config,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OpenAIParser) ParseFurniture(ctx context.Context, description string, name string) (domain.FurnitureDefinition, error) {
	if strings.TrimSpace(description) == "" {
		return domain.FurnitureDefinition{}, ErrEmptyDescription
	}
	if p.config.APIKey == "" {
		return domain.FurnitureDefinition{}, ErrOpenAINotConfigured
	}

	userPrompt := fmt.Sprintf("Description:\n%s\n\nSuggested name: %s\n\nReturn FurnitureDefinition JSON only.", description, name)
	raw, err := p.chatCompletion(ctx, userPrompt)
	if err != nil {
		return domain.FurnitureDefinition{}, err
	}

	return parseAndValidateFurnitureJSON(raw)
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
	Temperature    float64         `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (p *OpenAIParser) chatCompletion(ctx context.Context, userPrompt string) ([]byte, error) {
	payload := chatRequest{
		Model: p.config.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		ResponseFormat: &responseFormat{Type: "json_object"},
		Temperature:    0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(p.config.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, ErrOpenAIRequest
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: %s", ErrOpenAIRequest, string(respBody))
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("%w: %s", ErrOpenAIRequest, parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return nil, ErrOpenAIRequest
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return []byte(strings.TrimSpace(content)), nil
}
