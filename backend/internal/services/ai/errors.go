package ai

import "errors"

var (
	ErrEmptyDescription    = errors.New("description is required")
	ErrInvalidAIOutput     = errors.New("ai output is not valid furniture definition")
	ErrOpenAINotConfigured = errors.New("openai api key is not configured")
	ErrOpenAIRequest       = errors.New("openai request failed")
)
