package ai

import (
	"strings"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
)

func TestRepairMisalignedChildConstraints(t *testing.T) {
	raw := []byte(`{
		"id": "",
		"name": "",
		"description": "Closet empotrado 2.40m, dos cuerpos",
		"root": {
			"id": "root",
			"constraints": {
				"width": {"mode": "fill"},
				"height": {"mode": "fill"},
				"depth": {"mode": "fixed", "value": 600}
			},
			"split": {"axis": "x", "ratios": [0.5, 0.5]},
			"children": [
				{
					"id": "left",
					"constraints": {
						"width": {"mode": "fill"},
						"height": {"mode": "fill"},
						"depth": {"mode": "fill"}
					},
					"features": ["shelf_set"],
					"fronts": ["door"]
				},
				{
					"id": "right",
					"constraints": {
						"width": {"mode": "fill"},
						"height": {"mode": "fill"},
						"depth": {"mode": "fill"}
					},
					"features": ["shelf_set"],
					"fronts": []
				}
			]
		}
	}`)

	furniture, err := parseAndValidateFurnitureJSON(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}

	if furniture.Name == "" || furniture.ID == "" {
		t.Fatal("expected repaired furniture metadata")
	}
	if furniture.Root.Children[0].Constraints.Width == nil ||
		furniture.Root.Children[0].Constraints.Width.Mode != "ratio" ||
		furniture.Root.Children[0].Constraints.Width.Value != 0.5 {
		t.Fatalf("left width constraint not repaired: %+v", furniture.Root.Children[0].Constraints.Width)
	}
}

func TestRepairChildrenWithoutSplit(t *testing.T) {
	raw := []byte(`{
		"id": "closet-1",
		"name": "Closet",
		"description": "demo",
		"root": {
			"id": "root",
			"constraints": {
				"width": {"mode": "fill"},
				"height": {"mode": "fill"},
				"depth": {"mode": "fill"}
			},
			"children": [
				{"id": "a", "constraints": {"width": {"mode": "fill"}, "height": {"mode": "fill"}, "depth": {"mode": "fill"}}},
				{"id": "b", "constraints": {"width": {"mode": "fill"}, "height": {"mode": "fill"}, "depth": {"mode": "fill"}}}
			]
		}
	}`)

	furniture, err := parseAndValidateFurnitureJSON(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}
	if furniture.Root.Split == nil {
		t.Fatal("expected auto-created split")
	}
}

func TestInvalidOutputIncludesValidationSummary(t *testing.T) {
	summary := summarizeValidationErrors(volumetree.ValidationResult{
		Valid: false,
		Errors: []volumetree.ValidationError{
			{Field: "root.split.ratios", Message: "ratios must sum to 1"},
		},
	})
	if !strings.Contains(summary, "root.split.ratios") {
		t.Fatalf("summary = %q", summary)
	}
}
