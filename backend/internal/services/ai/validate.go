package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
)

func parseAndValidateFurnitureJSON(raw []byte) (domain.FurnitureDefinition, error) {
	normalized, err := normalizeAIJSON(raw)
	if err != nil {
		return domain.FurnitureDefinition{}, err
	}

	var furniture domain.FurnitureDefinition
	if err := json.Unmarshal(normalized, &furniture); err != nil {
		return domain.FurnitureDefinition{}, err
	}

	if result := volumetree.ValidateFurnitureDefinition(furniture); !result.Valid {
		return domain.FurnitureDefinition{}, fmt.Errorf("%w: %s", ErrInvalidAIOutput, summarizeValidationErrors(result))
	}

	return furniture, nil
}

func summarizeValidationErrors(result volumetree.ValidationResult) string {
	if len(result.Errors) == 0 {
		return "unknown validation failure"
	}

	limit := len(result.Errors)
	if limit > 3 {
		limit = 3
	}

	parts := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		err := result.Errors[i]
		parts = append(parts, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	if len(result.Errors) > limit {
		parts = append(parts, fmt.Sprintf("+%d more", len(result.Errors)-limit))
	}
	return strings.Join(parts, "; ")
}
