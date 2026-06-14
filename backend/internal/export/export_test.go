package export_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/costing"
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/export"
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
		t.Fatalf("read fixture: %v", err)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return value
}

func compileExampleCloset(t *testing.T) (domain.ManufacturingModel, domain.CostResult) {
	t.Helper()
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	furniture := loadFixture[domain.FurnitureDefinition](t, "example-closet.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	model, err := manufacturing.CompileManufacturing(resolved)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	cost, err := costing.CalculateCost(model)
	if err != nil {
		t.Fatalf("cost: %v", err)
	}
	return model, cost
}

func TestBuildBOM(t *testing.T) {
	model, cost := compileExampleCloset(t)
	bom, err := export.BuildBOM(model, export.BuildOptions{
		FurnitureName: "Closet demo",
		Cost:          &cost,
	})
	if err != nil {
		t.Fatalf("build bom: %v", err)
	}
	if len(bom.Parts) == 0 {
		t.Fatal("expected parts in BOM")
	}
	if bom.Summary.PartCount != len(model.Parts) {
		t.Fatalf("part count = %d, want %d", bom.Summary.PartCount, len(model.Parts))
	}
	if bom.Cost == nil {
		t.Fatal("expected cost in BOM")
	}
}

func TestBuildCutPlan(t *testing.T) {
	model, _ := compileExampleCloset(t)
	plan, err := export.BuildCutPlan(model, export.BuildOptions{FurnitureName: "Closet demo"})
	if err != nil {
		t.Fatalf("build cut plan: %v", err)
	}
	if len(plan.Sheets) == 0 {
		t.Fatal("expected cut sheets")
	}
	if len(plan.Sheets[0].Parts) == 0 {
		t.Fatal("expected cut parts")
	}
}

func TestGeneratePDF(t *testing.T) {
	model, cost := compileExampleCloset(t)
	opts := export.BuildOptions{FurnitureName: "Closet demo", Cost: &cost}

	bom, err := export.BuildBOM(model, opts)
	if err != nil {
		t.Fatalf("bom: %v", err)
	}
	plan, err := export.BuildCutPlan(model, opts)
	if err != nil {
		t.Fatalf("cut plan: %v", err)
	}

	pdf, err := export.GeneratePDF(bom, plan)
	if err != nil {
		t.Fatalf("pdf: %v", err)
	}
	if len(pdf) < 1000 {
		t.Fatalf("pdf too small: %d bytes", len(pdf))
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF")) {
		t.Fatal("output is not a PDF")
	}
}
