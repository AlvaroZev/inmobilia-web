package cases

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
	"github.com/inmobilia/inmobilia-web/backend/internal/resolvedfurniture"
)

// Caso 01: un volumen hoja con drawer_stack count=1 (cajonera simple, sin escritorio).
func TestCase01SingleDrawer_SolveAndValidate(t *testing.T) {
	resolved := solveFixture(t, "example-single-drawer.json")

	if resolved.ID != "single-drawer-001" {
		t.Fatalf("id = %q, want single-drawer-001", resolved.ID)
	}

	result := resolvedfurniture.ValidateResolvedFurniture(resolved)
	if !result.Valid {
		t.Fatalf("invalid resolved furniture: %+v", result.Errors)
	}

	root := resolved.Root
	if len(root.Children) != 0 {
		t.Fatalf("expected leaf root, got %d children", len(root.Children))
	}
	if math.Abs(root.Width-500) > 1 {
		t.Fatalf("width = %v, want 500", root.Width)
	}
	if math.Abs(root.Height-350) > 1 {
		t.Fatalf("height = %v, want 350", root.Height)
	}
	if math.Abs(root.Depth-400) > 1 {
		t.Fatalf("depth = %v, want 400", root.Depth)
	}

	if len(root.Features) != 1 {
		t.Fatalf("feature count = %d, want 1", len(root.Features))
	}
	feature := root.Features[0]
	if feature.Type != "drawer_stack" {
		t.Fatalf("feature type = %q, want drawer_stack", feature.Type)
	}
	if !featureMatchesVolume(feature, root) {
		t.Fatalf("drawer_stack bbox does not match root volume")
	}

	var params map[string]any
	if err := json.Unmarshal(feature.Params, &params); err != nil {
		t.Fatalf("unmarshal feature params: %v", err)
	}
	if int(params["count"].(float64)) != 1 {
		t.Fatalf("drawer count = %v, want 1", params["count"])
	}

	if len(root.Fronts) != 0 {
		t.Fatalf("front count = %d, want 0 (frente solo desde drawer_stack)", len(root.Fronts))
	}
}

func TestCase01SingleDrawer_Manufacturing(t *testing.T) {
	resolved := solveFixture(t, "example-single-drawer.json")
	model, err := manufacturing.CompileManufacturing(resolved)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

  // Referencia 500×350×400 con fondos nordex 3 mm: frente 492×342, laterales 390×267, fondo 413×366.
	var front domain.Part
	var sides []domain.Part
	var bottom domain.Part
	for _, part := range model.Parts {
		switch part.Type {
		case string(domain.PartDoor):
			if part.Name == "Cajón Cuerpo 1 frente" {
				front = part
			}
		case string(domain.PartDrawerSide):
			sides = append(sides, part)
		case string(domain.PartDrawerBottom):
			bottom = part
		}
	}
	if math.Abs(front.Width-492) > 1 || math.Abs(front.Height-342) > 1 {
		t.Fatalf("front = %.0f×%.0f, want 492×342", front.Width, front.Height)
	}
	if len(sides) != 2 {
		t.Fatalf("drawer side count = %d, want 2", len(sides))
	}
	for _, side := range sides {
		if math.Abs(side.Width-390) > 1 || math.Abs(side.Height-267) > 1 {
			t.Fatalf("side = %.0f×%.0f, want 390×267", side.Width, side.Height)
		}
	}
	if math.Abs(bottom.Width-413) > 1 || math.Abs(bottom.Height-366) > 1 {
		t.Fatalf("bottom = %.0f×%.0f, want 413×366", bottom.Width, bottom.Height)
	}
	if bottom.Material.Thickness != float64(manufacturing.NordexThicknessMm) {
		t.Fatalf("bottom thickness = %.0f, want nordex %d", bottom.Material.Thickness, manufacturing.NordexThicknessMm)
	}
	if countHardware(model.Hardware, "drawer_runner") != 1 {
		t.Fatalf("runner count = %d, want 1", countHardware(model.Hardware, "drawer_runner"))
	}
}

func TestCase01SingleDrawer_JSONRoundTrip(t *testing.T) {
	resolved := solveFixture(t, "example-single-drawer.json")

	data, err := json.Marshal(resolved)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded domain.ResolvedFurniture
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	result := resolvedfurniture.ValidateResolvedFurniture(decoded)
	if !result.Valid {
		t.Fatalf("round-trip invalid: %+v", result.Errors)
	}
}

func featureMatchesVolume(feature domain.ResolvedFeature, volume domain.ResolvedVolume) bool {
	return math.Abs(feature.X-volume.X) < 0.01 &&
		math.Abs(feature.Y-volume.Y) < 0.01 &&
		math.Abs(feature.Z-volume.Z) < 0.01 &&
		math.Abs(feature.Width-volume.Width) < 0.01 &&
		math.Abs(feature.Height-volume.Height) < 0.01 &&
		math.Abs(feature.Depth-volume.Depth) < 0.01
}

func countParts(parts []domain.Part, partType string) int {
	n := 0
	for _, part := range parts {
		if part.Type == partType {
			n++
		}
	}
	return n
}

func countHardware(items []domain.Hardware, hwType string) int {
	n := 0
	for _, item := range items {
		if item.Type == hwType {
			n++
		}
	}
	return n
}
