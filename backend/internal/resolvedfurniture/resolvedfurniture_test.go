package resolvedfurniture_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/resolvedfurniture"
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

func TestValidateSolverOutput(t *testing.T) {
	resolved := solveExampleCloset(t)
	result := resolvedfurniture.ValidateResolvedFurniture(resolved)
	if !result.Valid {
		t.Fatalf("expected valid resolved furniture, got: %+v", result.Errors)
	}
}

func TestResolvedTreeMetrics(t *testing.T) {
	resolved := solveExampleCloset(t)
	root := resolved.Root

	if resolvedfurniture.GetNodeCount(root) != 5 {
		t.Fatalf("node count = %d, want 5", resolvedfurniture.GetNodeCount(root))
	}
	if resolvedfurniture.GetTreeDepth(root) != 2 {
		t.Fatalf("tree depth = %d, want 2", resolvedfurniture.GetTreeDepth(root))
	}
	if len(resolvedfurniture.GetLeafVolumes(root)) != 3 {
		t.Fatalf("leaf count = %d, want 3", len(resolvedfurniture.GetLeafVolumes(root)))
	}
}

func TestChildContainment(t *testing.T) {
	resolved := solveExampleCloset(t)

	for _, ref := range resolvedfurniture.FlattenResolvedTree(resolved.Root) {
		if ref.Parent == nil {
			continue
		}
		if !resolvedfurniture.VolumeContains(*ref.Parent, ref.Volume) {
			t.Fatalf("volume %s is not contained in parent %s", ref.Volume.ID, ref.Parent.ID)
		}
	}
}

func TestExternalDimensions(t *testing.T) {
	resolved := solveExampleCloset(t)
	width, height, depth := resolvedfurniture.ExternalDimensions(resolved)

	if math.Abs(width-4098) > 1 {
		t.Fatalf("width = %v, want 4098", width)
	}
	if math.Abs(height-2548) > 1 {
		t.Fatalf("height = %v, want 2548", height)
	}
	if math.Abs(depth-600) > 1 {
		t.Fatalf("depth = %v, want 600", depth)
	}
}

func TestValidateRejectsNegativeDimension(t *testing.T) {
	resolved := solveExampleCloset(t)
	resolved.Root.Width = -10

	result := resolvedfurniture.ValidateResolvedFurniture(resolved)
	if result.Valid {
		t.Fatal("expected invalid resolved furniture for negative width")
	}
}

func TestValidateRejectsOverlappingSiblings(t *testing.T) {
	resolved := solveExampleCloset(t)
	resolved.Root.Children[0].X = resolved.Root.Children[1].X

	result := resolvedfurniture.ValidateResolvedFurniture(resolved)
	if result.Valid {
		t.Fatal("expected invalid resolved furniture for overlapping siblings")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	resolved := solveExampleCloset(t)

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
		t.Fatalf("round-trip output invalid: %+v", result.Errors)
	}
}
