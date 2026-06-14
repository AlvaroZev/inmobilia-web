package ai

import (
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
)

func TestNormalizeAIJSONStringFeatures(t *testing.T) {
	raw := []byte(`{
		"id": "test",
		"name": "Closet",
		"description": "demo",
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
						"width": {"mode": "ratio", "value": 0.5},
						"height": {"mode": "fill"},
						"depth": {"mode": "fill"}
					},
					"features": ["shelf_set"],
					"fronts": ["door"]
				},
				{
					"id": "right",
					"constraints": {
						"width": {"mode": "ratio", "value": 0.5},
						"height": {"mode": "fill"},
						"depth": {"mode": "fill"}
					},
					"features": ["shelf_set", "hanger_rod"],
					"fronts": []
				}
			],
			"features": [],
			"fronts": []
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

	if len(furniture.Root.Children[0].Features) != 1 {
		t.Fatalf("left features = %d, want 1", len(furniture.Root.Children[0].Features))
	}
	if furniture.Root.Children[0].Features[0].Type != "shelf_set" {
		t.Fatalf("feature type = %q", furniture.Root.Children[0].Features[0].Type)
	}
	if len(furniture.Root.Children[1].Features) != 2 {
		t.Fatalf("right features = %d, want 2", len(furniture.Root.Children[1].Features))
	}
}
