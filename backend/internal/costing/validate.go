package costing

import (
	"fmt"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

func (r *ValidationResult) add(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
	r.Valid = false
}

func ValidateCostResult(result domain.CostResult) ValidationResult {
	validation := ValidationResult{Valid: true, Errors: []ValidationError{}}

	if result.FurnitureID == "" {
		validation.add("furnitureId", "furnitureId is required")
	}
	if result.Currency == "" {
		validation.add("currency", "currency is required")
	}
	if result.Total < 0 {
		validation.add("total", "total must be non-negative")
	}

	expectedSubtotal := sumMaterialTotals(result.Materials) +
		sumHardwareTotals(result.Hardware) +
		result.Labor.Total

	if diff := result.Subtotal - expectedSubtotal; diff > 0.02 || diff < -0.02 {
		validation.add("subtotal", fmt.Sprintf("subtotal %v does not match line items %v", result.Subtotal, expectedSubtotal))
	}

	expectedTotal := result.Subtotal + result.Waste.Total
	if diff := result.Total - expectedTotal; diff > 0.02 || diff < -0.02 {
		validation.add("total", fmt.Sprintf("total %v does not match subtotal + waste %v", result.Total, expectedTotal))
	}

	return validation
}
