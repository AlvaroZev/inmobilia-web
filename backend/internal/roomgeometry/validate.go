package roomgeometry

import (
	"fmt"
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

func (r *ValidationResult) add(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
	r.Valid = false
}

func ValidateRoomGeometry(room domain.RoomGeometry) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	if room.ID == "" {
		result.add("id", "id is required")
	}

	validatePolygon(room.Perimeter, "perimeter", &result)
	validatePlane(room.Floor, "floor", &result)
	validatePlane(room.Ceiling, "ceiling", &result)

	ids := map[string]string{}
	trackID := func(id, field string) {
		if id == "" {
			return
		}
		if existing, ok := ids[id]; ok {
			result.add(field, fmt.Sprintf(`duplicate id "%s" (also used in %s)`, id, existing))
			return
		}
		ids[id] = field
	}

	for i, wall := range room.Walls {
		field := fmt.Sprintf("walls[%d].id", i)
		trackID(wall.ID, field)
		validateWall(wall, fmt.Sprintf("walls[%d]", i), &result)
	}

	wallIndex := map[string]domain.Wall{}
	for _, wall := range room.Walls {
		wallIndex[wall.ID] = wall
	}

	for i, opening := range room.Openings {
		field := fmt.Sprintf("openings[%d].id", i)
		trackID(opening.ID, field)
		wall, ok := wallIndex[opening.WallID]
		if !ok {
			validateOpening(opening, nil, fmt.Sprintf("openings[%d]", i), &result)
			continue
		}
		validateOpening(opening, &wall, fmt.Sprintf("openings[%d]", i), &result)
	}

	for i, obstacle := range room.Obstacles {
		field := fmt.Sprintf("obstacles[%d].id", i)
		trackID(obstacle.ID, field)
		validateObstacle(obstacle, fmt.Sprintf("obstacles[%d]", i), &result)
	}

	return result
}

func validatePolygon(polygon domain.Polygon2D, prefix string, result *ValidationResult) {
	if len(polygon.Vertices) < 3 {
		result.add(prefix+".vertices", "at least 3 vertices required")
		return
	}

	for i := range polygon.Vertices {
		j := (i + 1) % len(polygon.Vertices)
		a := polygon.Vertices[i]
		b := polygon.Vertices[j]
		if math.Hypot(a.X-b.X, a.Y-b.Y) < geom.Epsilon {
			result.add(fmt.Sprintf("%s.vertices[%d]", prefix, i), "zero-length edge")
		}
	}

	if math.Abs(geom.PolygonArea2D(polygon)) < geom.Epsilon {
		result.add(prefix, "polygon area must be greater than zero")
	}
}

func validatePlane(plane domain.Plane, field string, result *ValidationResult) {
	if geom.Magnitude3(plane.Normal) < geom.Epsilon {
		result.add(field+".normal", "normal vector must be non-zero")
	}
}

func validateWall(wall domain.Wall, prefix string, result *ValidationResult) {
	if wall.ID == "" {
		result.add(prefix+".id", "id is required")
	}
	if wall.Thickness <= geom.Epsilon {
		result.add(prefix+".thickness", "thickness must be greater than zero")
	}
	if len(wall.Vertices) < 3 {
		result.add(prefix+".vertices", "at least 3 vertices required")
		return
	}

	if _, ok := GetWallPlane(wall); !ok {
		result.add(prefix+".vertices", "vertices must define a non-degenerate face")
	}

	if len(wall.Vertices) >= 4 {
		frame, ok := GetWallLocalFrame(wall)
		if !ok {
			result.add(prefix+".vertices", "quad wall frame is degenerate")
		} else if frame.Width < geom.Epsilon || frame.Height < geom.Epsilon {
			result.add(prefix+".vertices", "wall width and height must be greater than zero")
		}
	}
}

func validateOpening(opening domain.Opening, wall *domain.Wall, prefix string, result *ValidationResult) {
	if opening.ID == "" {
		result.add(prefix+".id", "id is required")
	}
	if opening.Width <= geom.Epsilon {
		result.add(prefix+".width", "width must be greater than zero")
	}
	if opening.Height <= geom.Epsilon {
		result.add(prefix+".height", "height must be greater than zero")
	}
	if wall == nil {
		result.add(prefix+".wallId", fmt.Sprintf(`wall "%s" not found`, opening.WallID))
		return
	}

	frame, ok := GetWallLocalFrame(*wall)
	if !ok {
		result.add(prefix+".wallId", "wall geometry is not a valid quad face")
		return
	}

	dist := math.Abs(geom.SignedDistanceToPlane(opening.Origin, domain.Plane{
		Point:  frame.Origin,
		Normal: frame.Normal,
	}))
	if dist > 5 {
		result.add(prefix+".origin", "opening origin is not on the wall face")
	}

	u, v, okLocal := ProjectPointOnWallLocal(*wall, opening.Origin)
	if !okLocal {
		return
	}

	if u < -geom.Epsilon || v < -geom.Epsilon {
		result.add(prefix+".origin", "opening origin is outside wall bounds")
	}
	if u+opening.Width > frame.Width+geom.Epsilon {
		result.add(prefix+".width", "opening exceeds wall width")
	}
	if v+opening.Height > frame.Height+geom.Epsilon {
		result.add(prefix+".height", "opening exceeds wall height")
	}
}

func validateObstacle(obstacle domain.Obstacle, prefix string, result *ValidationResult) {
	if obstacle.ID == "" {
		result.add(prefix+".id", "id is required")
	}

	bounds := geom.NormalizeBounds(obstacle.Bounds)
	if geom.BoundsVolume(bounds) < geom.Epsilon {
		result.add(prefix+".bounds", "bounds must have positive volume")
	}
}
