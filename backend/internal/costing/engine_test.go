package costing_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/costing"
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

func compileExampleCloset(t *testing.T) domain.ManufacturingModel {
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
	return model
}

func TestCalculateExampleClosetCost(t *testing.T) {
	model := compileExampleCloset(t)
	result, err := costing.CalculateCost(model)
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}

	if result.FurnitureID != "closet-001" {
		t.Fatalf("furnitureId = %q, want closet-001", result.FurnitureID)
	}
	if result.Currency != "USD" {
		t.Fatalf("currency = %q, want USD", result.Currency)
	}
	if len(result.Materials) == 0 {
		t.Fatal("expected material lines")
	}
	if len(result.Hardware) == 0 {
		t.Fatal("expected hardware lines")
	}
	if result.Labor.Hours <= 0 {
		t.Fatal("expected positive labor hours")
	}
	if result.Waste.Total <= 0 {
		t.Fatal("expected positive waste cost")
	}
	if result.Total <= result.Subtotal {
		t.Fatal("total should include waste on top of subtotal")
	}

	validation := costing.ValidateCostResult(result)
	if !validation.Valid {
		t.Fatalf("invalid cost result: %+v", validation.Errors)
	}
}

func TestCostTotalsAreConsistent(t *testing.T) {
	model := compileExampleCloset(t)
	result, err := costing.CalculateCost(model)
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}

	var materialsTotal float64
	for _, line := range result.Materials {
		materialsTotal += line.Total
	}
	var hardwareTotal float64
	for _, line := range result.Hardware {
		hardwareTotal += line.Total
	}

	expectedSubtotal := materialsTotal + hardwareTotal + result.Labor.Total
	if math.Abs(result.Subtotal-expectedSubtotal) > 0.05 {
		t.Fatalf("subtotal = %v, want %v", result.Subtotal, expectedSubtotal)
	}

	expectedTotal := result.Subtotal + result.Waste.Total
	if math.Abs(result.Total-expectedTotal) > 0.05 {
		t.Fatalf("total = %v, want %v", result.Total, expectedTotal)
	}
}

func TestCalculateRejectsInvalidModel(t *testing.T) {
	_, err := costing.CalculateCost(domain.ManufacturingModel{})
	if err != costing.ErrInvalidManufacturingModel {
		t.Fatalf("expected ErrInvalidManufacturingModel, got %v", err)
	}
}

func TestCostJSONRoundTrip(t *testing.T) {
	model := compileExampleCloset(t)
	result, err := costing.CalculateCost(model)
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded domain.CostResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	validation := costing.ValidateCostResult(decoded)
	if !validation.Valid {
		t.Fatalf("round-trip invalid: %+v", validation.Errors)
	}
}

func TestCustomCostOptions(t *testing.T) {
	model := compileExampleCloset(t)
	opts := costing.DefaultCostOptions()
	opts.Currency = "COP"
	opts.LaborRatePerHour = 50000

	result, err := costing.CalculateCostWithOptions(model, opts)
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}
	if result.Currency != "COP" {
		t.Fatalf("currency = %q, want COP", result.Currency)
	}
	if result.Labor.RatePerHour != 50000 {
		t.Fatalf("labor rate = %v, want 50000", result.Labor.RatePerHour)
	}
}
