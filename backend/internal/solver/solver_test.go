package solver_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/solver"
	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
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

func TestSolveSimpleFillBox(t *testing.T) {
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")

	furniture := domain.FurnitureDefinition{
		ID:   "simple-001",
		Name: "Caja simple",
		Root: domain.VolumeNode{
			ID: "root",
			Constraints: domain.VolumeConstraints{
				Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
			},
			Children: []domain.VolumeNode{},
			Features: []domain.Feature{},
			Fronts:   []domain.Front{},
		},
	}

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}

	root := resolved.Root
	if math.Abs(root.X-50) > 1 {
		t.Fatalf("x = %v, want 50", root.X)
	}
	if math.Abs(root.Y-100) > 1 {
		t.Fatalf("y = %v, want 100 (floorOffset)", root.Y)
	}
	if math.Abs(root.Z-10) > 1 {
		t.Fatalf("z = %v, want 10", root.Z)
	}
	if math.Abs(root.Width-4098) > 1 {
		t.Fatalf("width = %v, want 4098", root.Width)
	}
	if math.Abs(root.Height-2548) > 1 {
		t.Fatalf("height = %v, want 2548", root.Height)
	}
	if math.Abs(root.Depth-609) > 1 {
		t.Fatalf("depth = %v, want 609", root.Depth)
	}
}

func TestSolveExampleCloset(t *testing.T) {
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	furniture := loadFixture[domain.FurnitureDefinition](t, "example-closet.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}

	root := resolved.Root
	if len(root.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(root.Children))
	}
	if math.Abs(root.Depth-600) > 1 {
		t.Fatalf("root depth = %v, want 600 (fixed depth)", root.Depth)
	}

	left := root.Children[0]
	right := root.Children[1]
	if math.Abs(left.Width-root.Width/2) > 1 {
		t.Fatalf("left width = %v, want half of root %v", left.Width, root.Width)
	}
	if math.Abs(right.Width-root.Width/2) > 1 {
		t.Fatalf("right width = %v, want half of root", right.Width)
	}

	if len(right.Children) != 2 {
		t.Fatalf("right children = %d, want 2", len(right.Children))
	}
	upper := right.Children[0]
	drawers := right.Children[1]
	if math.Abs(upper.Height-right.Height*0.7) > 2 {
		t.Fatalf("upper height = %v, want ~70%% of right", upper.Height)
	}
	if math.Abs(drawers.Height-right.Height*0.3) > 2 {
		t.Fatalf("drawers height = %v, want ~30%% of right", drawers.Height)
	}

	if len(volumetree.CollectFeatures(furniture.Root)) != len(collectResolvedFeatures(resolved.Root)) {
		t.Fatal("feature count mismatch between definition and resolved")
	}
	if resolved.Root.MaterialID != "melamine-white-18" {
		t.Fatalf("materialId = %q, want melamine-white-18", resolved.Root.MaterialID)
	}
}

func TestSolveRejectsInvalidFurniture(t *testing.T) {
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")
	furniture := loadFixture[domain.FurnitureDefinition](t, "example-closet.json")
	furniture.Root.Split.Ratios = []float64{0.6, 0.6}

	_, err := solver.SolveConstraints(room, furniture, installation)
	if err != solver.ErrInvalidFurniture {
		t.Fatalf("expected ErrInvalidFurniture, got %v", err)
	}
}

func TestSolveRejectsOversizedFixedDepth(t *testing.T) {
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")

	furniture := domain.FurnitureDefinition{
		ID:   "deep-001",
		Name: "Demasiado profundo",
		Root: domain.VolumeNode{
			ID: "root",
			Constraints: domain.VolumeConstraints{
				Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFixed, Value: 2000},
			},
			Children: []domain.VolumeNode{},
		},
	}

	_, err := solver.SolveConstraints(room, furniture, installation)
	if err != solver.ErrDimensionExceedsSpace {
		t.Fatalf("expected ErrDimensionExceedsSpace, got %v", err)
	}
}

func collectResolvedFeatures(volume domain.ResolvedVolume) []domain.ResolvedFeature {
	features := append([]domain.ResolvedFeature{}, volume.Features...)
	for _, child := range volume.Children {
		features = append(features, collectResolvedFeatures(child)...)
	}
	return features
}
