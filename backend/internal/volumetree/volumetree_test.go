package volumetree_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
)

func loadExampleCloset(t *testing.T) domain.FurnitureDefinition {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	fixturePath := filepath.Join(
		filepath.Dir(file),
		"..", "..", "..", "frontend", "src", "domain", "fixtures", "example-closet.json",
	)

	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var furniture domain.FurnitureDefinition
	if err := json.Unmarshal(data, &furniture); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	return furniture
}

func TestValidateExampleCloset(t *testing.T) {
	furniture := loadExampleCloset(t)
	result := volumetree.ValidateFurnitureDefinition(furniture)

	if !result.Valid {
		t.Fatalf("expected valid furniture, got errors: %+v", result.Errors)
	}
}

func TestTreeTraversal(t *testing.T) {
	furniture := loadExampleCloset(t)
	root := furniture.Root

	if volumetree.GetNodeCount(root) != 5 {
		t.Fatalf("node count = %d, want 5", volumetree.GetNodeCount(root))
	}
	if volumetree.GetTreeDepth(root) != 2 {
		t.Fatalf("tree depth = %d, want 2", volumetree.GetTreeDepth(root))
	}
	if len(volumetree.GetLeafNodes(root)) != 3 {
		t.Fatalf("leaf count = %d, want 3", len(volumetree.GetLeafNodes(root)))
	}

	ref, ok := volumetree.FindVolumeNodeByID(root, "right-drawers")
	if !ok {
		t.Fatal("expected to find right-drawers node")
	}
	if ref.Depth != 2 {
		t.Fatalf("right-drawers depth = %d, want 2", ref.Depth)
	}
	if ref.Parent == nil || ref.Parent.ID != "right-body" {
		t.Fatalf("expected parent right-body, got %+v", ref.Parent)
	}
}

func TestFeatureCollection(t *testing.T) {
	furniture := loadExampleCloset(t)
	counts := volumetree.CountFeaturesByType(furniture.Root)

	if counts["shelf_set"] != 2 {
		t.Fatalf("shelf_set count = %d, want 2", counts["shelf_set"])
	}
	if counts["drawer_stack"] != 1 {
		t.Fatalf("drawer_stack count = %d, want 1", counts["drawer_stack"])
	}
	if counts["hanger_rod"] != 1 {
		t.Fatalf("hanger_rod count = %d, want 1", counts["hanger_rod"])
	}
	if len(volumetree.CollectFronts(furniture.Root)) != 3 {
		t.Fatalf("front count = %d, want 3", len(volumetree.CollectFronts(furniture.Root)))
	}
}

func TestConstraintSummary(t *testing.T) {
	furniture := loadExampleCloset(t)
	summary := volumetree.SummarizeConstraints(furniture.Root)

	if !summary.HasFill {
		t.Fatal("expected root to have fill constraints")
	}
	if !summary.HasFixed {
		t.Fatal("expected root to have fixed depth constraint")
	}
	if summary.Depth == nil || summary.Depth.Mode != domain.DimensionFixed {
		t.Fatalf("expected fixed depth, got %+v", summary.Depth)
	}
	if math.Abs(summary.Depth.Value-600) > 1 {
		t.Fatalf("depth = %v, want 600", summary.Depth.Value)
	}
}

func TestValidateRejectsBadSplitRatios(t *testing.T) {
	furniture := loadExampleCloset(t)
	furniture.Root.Split.Ratios = []float64{0.6, 0.6}

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if result.Valid {
		t.Fatal("expected invalid furniture when ratios do not sum to 1")
	}
}

func TestValidateRejectsChildrenWithoutSplit(t *testing.T) {
	furniture := loadExampleCloset(t)
	furniture.Root.Split = nil

	result := volumetree.ValidateFurnitureDefinition(furniture)
	if result.Valid {
		t.Fatal("expected invalid furniture when children exist without split")
	}
}
