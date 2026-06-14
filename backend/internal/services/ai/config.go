package ai

import "os"

func provider() string {
	if value := os.Getenv("AI_PROVIDER"); value != "" {
		return value
	}
	if apiKey() != "" {
		return "openai"
	}
	return "mock"
}

func apiKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func model() string {
	if value := os.Getenv("OPENAI_MODEL"); value != "" {
		return value
	}
	return "gpt-4o-mini"
}

func baseURL() string {
	if value := os.Getenv("OPENAI_BASE_URL"); value != "" {
		return value
	}
	return "https://api.openai.com/v1"
}
