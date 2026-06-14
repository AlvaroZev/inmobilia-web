package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/api"
	"github.com/inmobilia/inmobilia-web/backend/internal/costing"
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
	"github.com/inmobilia/inmobilia-web/backend/internal/solver"
)

func TestExportBOMEndpoint(t *testing.T) {
	model, cost := compileForExport(t)
	router := api.NewRouter(ai.NewMockParser(), api.ServiceInfo{AIProvider: "mock"})

	body, _ := json.Marshal(api.ExportRequest{
		FurnitureName: "Closet",
		Model:         model,
		Cost:          &cost,
	})
	req := httptest.NewRequest(http.MethodPost, "/export/bom", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.Code, resp.Body.String())
	}
}

func TestExportPDFEndpoint(t *testing.T) {
	model, cost := compileForExport(t)
	router := api.NewRouter(ai.NewMockParser(), api.ServiceInfo{AIProvider: "mock"})

	body, _ := json.Marshal(api.ExportRequest{
		FurnitureName: "Closet",
		Model:         model,
		Cost:          &cost,
	})
	req := httptest.NewRequest(http.MethodPost, "/export/pdf", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d", resp.Code)
	}
	if ct := resp.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("content-type = %q", ct)
	}
}

func loadExportFixture[T any](t *testing.T, name string) T {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	path := filepath.Join(filepath.Dir(file), "..", "..", "..", "frontend", "src", "domain", "fixtures", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return value
}

func compileForExport(t *testing.T) (domain.ManufacturingModel, domain.CostResult) {
	t.Helper()
	room := loadExportFixture[domain.RoomGeometry](t, "example-room.json")
	furniture := loadExportFixture[domain.FurnitureDefinition](t, "example-closet.json")
	installation := loadExportFixture[domain.InstallationConstraints](t, "example-installation.json")

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
