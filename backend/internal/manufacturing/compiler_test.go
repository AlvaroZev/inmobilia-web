package manufacturing_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
	"github.com/inmobilia/inmobilia-web/backend/internal/solver"
)

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "frontend", "src", "domain", "fixtures", name)
}

func loadFixture[T any](t *testing.T, name string) T {
	t.Helper()
	data, err := os.ReadFile(fixturePath(t, name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("unmarshal fixture %s: %v", name, err)
	}
	return value
}

func solveExampleCloset(t *testing.T) domain.ResolvedFurniture {
	t.Helper()
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	furniture := loadFixture[domain.FurnitureDefinition](t, "example-closet.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	return resolved
}

func countPartsByType(parts []domain.Part, partType string) int {
	count := 0
	for _, part := range parts {
		if part.Type == partType {
			count++
		}
	}
	return count
}

func TestCompileExampleCloset(t *testing.T) {
	resolved := solveExampleCloset(t)
	model, err := manufacturing.CompileManufacturing(resolved)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	if model.FurnitureID != "closet-001" {
		t.Fatalf("furnitureId = %q, want closet-001", model.FurnitureID)
	}
	if len(model.Parts) < 15 {
		t.Fatalf("parts count = %d, want at least 15", len(model.Parts))
	}

	if countPartsByType(model.Parts, string(domain.PartLateral)) != 2 {
		t.Fatalf("lateral count = %d, want 2", countPartsByType(model.Parts, string(domain.PartLateral)))
	}
	if countPartsByType(model.Parts, string(domain.PartShelf)) != 7 {
		t.Fatalf("shelf count = %d, want 7", countPartsByType(model.Parts, string(domain.PartShelf)))
	}
	if countPartsByType(model.Parts, string(domain.PartDivider)) != 2 {
		t.Fatalf("divider count = %d, want 2", countPartsByType(model.Parts, string(domain.PartDivider)))
	}
	if countPartsByType(model.Parts, string(domain.PartDoor)) < 5 {
		t.Fatalf("door/front count = %d, want at least 5", countPartsByType(model.Parts, string(domain.PartDoor)))
	}

	if len(model.EdgeBanding) == 0 {
		t.Fatal("expected edge banding entries")
	}
	if len(model.Hardware) < 5 {
		t.Fatalf("hardware count = %d, want at least 5", len(model.Hardware))
	}

	result := manufacturing.ValidateManufacturingModel(model)
	if !result.Valid {
		t.Fatalf("invalid model: %+v", result.Errors)
	}
}

func TestShelfDimensions(t *testing.T) {
	resolved := solveExampleCloset(t)
	model, err := manufacturing.CompileManufacturing(resolved)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	var shelf domain.Part
	for _, part := range model.Parts {
		if part.Type == string(domain.PartShelf) && part.VolumeID == "left-body" {
			shelf = part
			break
		}
	}
	if shelf.ID == "" {
		t.Fatal("expected shelf for left-body")
	}

	expectedWidth := resolved.Root.Children[0].Width - 36
	if math.Abs(shelf.Width-expectedWidth) > 1 {
		t.Fatalf("shelf width = %v, want ~%v", shelf.Width, expectedWidth)
	}
}

func TestCompileRejectsInvalidResolved(t *testing.T) {
	resolved := solveExampleCloset(t)
	resolved.Root.Width = -1

	_, err := manufacturing.CompileManufacturing(resolved)
	if err != manufacturing.ErrInvalidResolvedFurniture {
		t.Fatalf("expected ErrInvalidResolvedFurniture, got %v", err)
	}
}

func TestManufacturingJSONRoundTrip(t *testing.T) {
	resolved := solveExampleCloset(t)
	model, err := manufacturing.CompileManufacturing(resolved)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	data, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded domain.ManufacturingModel
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	result := manufacturing.ValidateManufacturingModel(decoded)
	if !result.Valid {
		t.Fatalf("round-trip invalid: %+v", result.Errors)
	}
}
