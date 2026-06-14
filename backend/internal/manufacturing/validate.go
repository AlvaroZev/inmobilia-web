package manufacturing

import (
	"fmt"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
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

func ValidateManufacturingModel(model domain.ManufacturingModel) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	if model.FurnitureID == "" {
		result.add("furnitureId", "furnitureId is required")
	}
	if len(model.Parts) == 0 {
		result.add("parts", "at least one part is required")
	}

	partIDs := map[string]string{}
	for i, part := range model.Parts {
		prefix := fmt.Sprintf("parts[%d]", i)
		if part.ID == "" {
			result.add(prefix+".id", "id is required")
		} else if existing, ok := partIDs[part.ID]; ok {
			result.add(prefix+".id", fmt.Sprintf(`duplicate id "%s" (also used in %s)`, part.ID, existing))
		} else {
			partIDs[part.ID] = prefix
		}
		if part.Width <= geom.Epsilon || part.Height <= geom.Epsilon || part.Thickness <= geom.Epsilon {
			result.add(prefix, "part dimensions must be greater than zero")
		}
		if part.VolumeID == "" {
			result.add(prefix+".volumeId", "volumeId is required")
		}
	}

	for i, edge := range model.EdgeBanding {
		if edge.PartID == "" {
			result.add(fmt.Sprintf("edgeBanding[%d].partId", i), "partId is required")
		}
		if edge.Length <= geom.Epsilon {
			result.add(fmt.Sprintf("edgeBanding[%d].length", i), "length must be greater than zero")
		}
	}

	for i, hw := range model.Hardware {
		if hw.Type == "" {
			result.add(fmt.Sprintf("hardware[%d].type", i), "type is required")
		}
		if hw.Quantity <= 0 {
			result.add(fmt.Sprintf("hardware[%d].quantity", i), "quantity must be greater than zero")
		}
	}

	return result
}
