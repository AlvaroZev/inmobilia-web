package roomgeometry_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/roomgeometry"
)

func loadExampleRoom(t *testing.T) domain.RoomGeometry {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	fixturePath := filepath.Join(
		filepath.Dir(file),
		"..", "..", "..", "frontend", "src", "domain", "fixtures", "example-room.json",
	)

	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var room domain.RoomGeometry
	if err := json.Unmarshal(data, &room); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	return room
}

func TestValidateExampleRoom(t *testing.T) {
	room := loadExampleRoom(t)
	result := roomgeometry.ValidateRoomGeometry(room)

	if !result.Valid {
		t.Fatalf("expected valid room, got errors: %+v", result.Errors)
	}
}

func TestRoomComputedDimensions(t *testing.T) {
	room := loadExampleRoom(t)

	area := roomgeometry.GetRoomFloorArea(room)
	if math.Abs(area-15_120_000) > 1 {
		t.Fatalf("floor area = %v, want ~15120000 mm²", area)
	}

	height := roomgeometry.GetFloorCeilingHeight(room)
	if math.Abs(height-2700) > 1 {
		t.Fatalf("floor-ceiling height = %v, want 2700", height)
	}

	angles := roomgeometry.GetPerimeterInteriorAngles(room)
	for i, angle := range angles {
		if math.Abs(angle-math.Pi/2) > 0.01 {
			t.Fatalf("corner %d angle = %v, want ~π/2", i, angle)
		}
	}

	skirting := roomgeometry.GetSkirtingObstacles(room)
	if len(skirting) != 1 {
		t.Fatalf("skirting count = %d, want 1", len(skirting))
	}
	if math.Abs(roomgeometry.GetSkirtingHeight(skirting[0])-100) > 1 {
		t.Fatalf("skirting height = %v, want 100", roomgeometry.GetSkirtingHeight(skirting[0]))
	}
}

func TestWallLocalFrame(t *testing.T) {
	room := loadExampleRoom(t)
	wall := room.Walls[0]

	frame, ok := roomgeometry.GetWallLocalFrame(wall)
	if !ok {
		t.Fatal("expected valid wall frame")
	}
	if math.Abs(frame.Width-4200) > 1 {
		t.Fatalf("wall width = %v, want 4200", frame.Width)
	}
	if math.Abs(frame.Height-2700) > 1 {
		t.Fatalf("wall height = %v, want 2700", frame.Height)
	}
}

func TestValidateRejectsInvalidOpening(t *testing.T) {
	room := loadExampleRoom(t)
	room.Openings[0].Height = 5000

	result := roomgeometry.ValidateRoomGeometry(room)
	if result.Valid {
		t.Fatal("expected invalid room when opening exceeds wall height")
	}
}
