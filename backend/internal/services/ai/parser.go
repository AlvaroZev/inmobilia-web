package ai

import (
	"context"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

type Parser interface {
	ParseFurniture(ctx context.Context, description string, name string) (domain.FurnitureDefinition, error)
}

func NewParserFromEnv() Parser {
	parser, _ := ResolveParserFromEnv()
	return parser
}

func ResolveParserFromEnv() (Parser, string) {
	selected := provider()
	if selected == "openai" {
		key := apiKey()
		if key == "" {
			return NewMockParser(), "mock (AI_PROVIDER=openai but OPENAI_API_KEY is missing)"
		}
		return NewOpenAIParser(OpenAIConfig{
			APIKey:  key,
			Model:   model(),
			BaseURL: baseURL(),
		}), "openai (" + model() + ")"
	}
	return NewMockParser(), "mock"
}
