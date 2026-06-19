package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/solver"
)

func main() {
	furnitureFile := "example-single-drawer.json"
	outFile := "example-resolved-single-drawer.json"
	if len(os.Args) > 1 {
		furnitureFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		outFile = os.Args[2]
	}

	room := mustLoad[domain.RoomGeometry]("example-room.json")
	installation := mustLoad[domain.InstallationConstraints]("example-installation.json")
	furniture := mustLoad[domain.FurnitureDefinition](furnitureFile)

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		panic(err)
	}

	fixturesDir := filepath.Join("..", "..", "..", "frontend", "src", "domain", "fixtures")
	out := filepath.Join(fixturesDir, outFile)
	writeJSON(out, resolved)
	fmt.Println("wrote", out)
}

func mustLoad[T any](name string) T {
	path := filepath.Join("..", "..", "..", "frontend", "src", "domain", "fixtures", name)
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		panic(err)
	}
	return value
}

func writeJSON(path string, value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		panic(err)
	}
}
