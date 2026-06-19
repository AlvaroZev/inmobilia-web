package cases

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
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

func solveFixture(t *testing.T, furnitureFile string) domain.ResolvedFurniture {
	t.Helper()
	room := loadFixture[domain.RoomGeometry](t, "example-room.json")
	installation := loadFixture[domain.InstallationConstraints](t, "example-installation.json")
	furniture := loadFixture[domain.FurnitureDefinition](t, furnitureFile)

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		t.Fatalf("solve %s: %v", furnitureFile, err)
	}
	return resolved
}
