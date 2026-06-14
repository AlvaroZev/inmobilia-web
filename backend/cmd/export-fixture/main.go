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
	room := mustLoad[domain.RoomGeometry]("example-room.json")
	furniture := mustLoad[domain.FurnitureDefinition]("example-closet.json")
	installation := mustLoad[domain.InstallationConstraints]("example-installation.json")

	resolved, err := solver.SolveConstraints(room, furniture, installation)
	if err != nil {
		panic(err)
	}

	out := filepath.Join("..", "..", "..", "frontend", "src", "domain", "fixtures", "example-resolved-closet.json")
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
