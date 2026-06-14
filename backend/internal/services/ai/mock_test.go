package ai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
)

func TestMockParserCloset(t *testing.T) {
	parser := ai.NewMockParser()
	furniture, err := parser.ParseFurniture(
		context.Background(),
		"Closet empotrado 2.40m con dos cuerpos y cajonera inferior",
		"Closet demo",
	)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}
	if len(furniture.Root.Children) != 2 {
		t.Fatalf("children = %d, want 2", len(furniture.Root.Children))
	}
}

func TestMockParserThreeBodies(t *testing.T) {
	parser := ai.NewMockParser()
	furniture, err := parser.ParseFurniture(
		context.Background(),
		"Ropero empotrado de 3 cuerpos con cajonera",
		"Ropero 3 cuerpos",
	)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}
	if len(furniture.Root.Children) != 3 {
		t.Fatalf("children = %d, want 3", len(furniture.Root.Children))
	}
	if len(furniture.Root.Split.Ratios) != 3 {
		t.Fatalf("ratios = %d, want 3", len(furniture.Root.Split.Ratios))
	}
}

func TestMockParserDesk(t *testing.T) {
	parser := ai.NewMockParser()
	furniture, err := parser.ParseFurniture(
		context.Background(),
		"Escritorio melamina blanco con cajonera lateral 600mm",
		"Escritorio oficina",
	)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}
	if len(furniture.Root.Children) != 2 {
		t.Fatalf("children = %d, want 2", len(furniture.Root.Children))
	}
	if furniture.Root.Features[0].Type != "desk_frame" {
		t.Fatalf("expected desk_frame on root, got %s", furniture.Root.Features[0].Type)
	}
	if furniture.Root.Children[1].ID != "drawer-tower" {
		t.Fatalf("expected drawer-tower as second child, got %s", furniture.Root.Children[1].ID)
	}
}

func TestMockParserDeskSingleDrawer(t *testing.T) {
	parser := ai.NewMockParser()
	furniture, err := parser.ParseFurniture(
		context.Background(),
		"Escritorio melamina con un cajón lateral",
		"Escritorio 1 cajón",
	)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}
	if furniture.Root.Children[1].ID != "drawer-bay" {
		t.Fatalf("expected drawer-bay, got %s", furniture.Root.Children[1].ID)
	}
	if len(furniture.Root.Children[1].Features) != 1 {
		t.Fatalf("features = %d, want 1", len(furniture.Root.Children[1].Features))
	}
	params := string(furniture.Root.Children[1].Features[0].Params)
	if !strings.Contains(params, `"count":1`) {
		t.Fatalf("expected count 1, got %s", params)
	}
	if !strings.Contains(params, `"drawerMode":"single"`) {
		t.Fatalf("expected drawerMode single, got %s", params)
	}
	if !strings.Contains(params, `"hasBase":false`) {
		t.Fatalf("expected hasBase false, got %s", params)
	}
}

func TestMockParserEntertainmentCenter(t *testing.T) {
	parser := ai.NewMockParser()
	furniture, err := parser.ParseFurniture(
		context.Background(),
		"Centro de entretenimiento 2.4m con nicho para TV y dos módulos con puertas",
		"Centro TV",
	)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if !result.Valid {
		t.Fatalf("invalid furniture: %+v", result.Errors)
	}
	if furniture.Root.Split == nil || furniture.Root.Split.Axis != "y" {
		t.Fatalf("expected vertical split, got %+v", furniture.Root.Split)
	}
	if len(furniture.Root.Children) != 2 {
		t.Fatalf("children = %d, want 2", len(furniture.Root.Children))
	}
}

func TestMockParserRejectsEmptyDescription(t *testing.T) {
	parser := ai.NewMockParser()
	_, err := parser.ParseFurniture(context.Background(), "  ", "")
	if err != ai.ErrEmptyDescription {
		t.Fatalf("expected ErrEmptyDescription, got %v", err)
	}
}
