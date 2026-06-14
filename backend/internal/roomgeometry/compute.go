package roomgeometry

import (
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
)

type WallLocalFrame struct {
	Origin domain.Point3D
	U      domain.Point3D
	V      domain.Point3D
	Normal domain.Point3D
	Width  float64
	Height float64
}

func FindWallByID(room domain.RoomGeometry, wallID string) (*domain.Wall, bool) {
	for i := range room.Walls {
		if room.Walls[i].ID == wallID {
			return &room.Walls[i], true
		}
	}
	return nil, false
}

func GetWallPlane(wall domain.Wall) (domain.Plane, bool) {
	if len(wall.Vertices) < 3 {
		return domain.Plane{}, false
	}
	return geom.PlaneFromPoints(wall.Vertices[0], wall.Vertices[1], wall.Vertices[2])
}

func GetWallLocalFrame(wall domain.Wall) (WallLocalFrame, bool) {
	if len(wall.Vertices) < 4 {
		return WallLocalFrame{}, false
	}

	v0 := wall.Vertices[0]
	v1 := wall.Vertices[1]
	v3 := wall.Vertices[3]

	uRaw := geom.Subtract3(v1, v0)
	vRaw := geom.Subtract3(v3, v0)
	u, okU := geom.Normalize3(uRaw)
	v, okV := geom.Normalize3(vRaw)
	if !okU || !okV {
		return WallLocalFrame{}, false
	}

	normal, okN := geom.Normalize3(geom.Cross3(u, v))
	if !okN {
		return WallLocalFrame{}, false
	}

	return WallLocalFrame{
		Origin: v0,
		U:      u,
		V:      v,
		Normal: normal,
		Width:  geom.Magnitude3(uRaw),
		Height: geom.Magnitude3(vRaw),
	}, true
}

func ProjectPointOnWallLocal(wall domain.Wall, point domain.Point3D) (u, v float64, ok bool) {
	frame, okFrame := GetWallLocalFrame(wall)
	if !okFrame {
		return 0, 0, false
	}

	onPlane := geom.ProjectPointOnPlane(point, domain.Plane{
		Point:  frame.Origin,
		Normal: frame.Normal,
	})
	rel := geom.Subtract3(onPlane, frame.Origin)

	return geom.Dot3(rel, frame.U), geom.Dot3(rel, frame.V), true
}

func GetWallBottomEdgeLength(wall domain.Wall) float64 {
	if len(wall.Vertices) < 2 {
		return 0
	}
	return geom.Distance3(wall.Vertices[0], wall.Vertices[1])
}

func GetWallHeight(wall domain.Wall) float64 {
	if len(wall.Vertices) < 4 {
		return 0
	}
	return geom.Distance3(wall.Vertices[0], wall.Vertices[3])
}

func GetWallBottomEdgeAngle(wall domain.Wall) float64 {
	if len(wall.Vertices) < 2 {
		return 0
	}
	edge := geom.Subtract3(wall.Vertices[1], wall.Vertices[0])
	return math.Atan2(edge.Z, edge.X)
}

func GetWallOutOfPlumbVector(wall domain.Wall) domain.Point3D {
	if wall.OutOfPlumb != nil {
		return *wall.OutOfPlumb
	}
	if len(wall.Vertices) < 4 {
		return domain.Point3D{}
	}

	vertical := geom.Subtract3(wall.Vertices[3], wall.Vertices[0])
	trueUp := domain.Point3D{Y: 1}
	projectedLength := geom.Dot3(vertical, trueUp)
	plumb := domain.Point3D{Y: projectedLength}
	return geom.Subtract3(vertical, plumb)
}

func GetWallOutOfPlumbMagnitude(wall domain.Wall) float64 {
	return geom.Magnitude3(GetWallOutOfPlumbVector(wall))
}

func GetFloorCeilingHeight(room domain.RoomGeometry) float64 {
	floor, okFloor := geom.NormalizePlane(room.Floor)
	ceiling, okCeiling := geom.NormalizePlane(room.Ceiling)
	if !okFloor || !okCeiling {
		return 0
	}
	return math.Abs(geom.SignedDistanceToPlane(ceiling.Point, floor))
}

func GetRoomFloorArea(room domain.RoomGeometry) float64 {
	return math.Abs(geom.PolygonArea2D(room.Perimeter))
}

func GetRoomPerimeterLength(room domain.RoomGeometry) float64 {
	return geom.PolygonPerimeter2D(room.Perimeter)
}

func GetPerimeterInteriorAngles(room domain.RoomGeometry) []float64 {
	angles := make([]float64, len(room.Perimeter.Vertices))
	for i := range room.Perimeter.Vertices {
		angles[i] = geom.InteriorAngleAtVertex2D(room.Perimeter, i)
	}
	return angles
}

func GetSkirtingObstacles(room domain.RoomGeometry) []domain.Obstacle {
	var result []domain.Obstacle
	for _, o := range room.Obstacles {
		if o.Type == domain.ObstacleSkirting {
			result = append(result, o)
		}
	}
	return result
}

func GetSkirtingHeight(obstacle domain.Obstacle) float64 {
	bounds := geom.NormalizeBounds(obstacle.Bounds)
	return bounds.Max.Y - bounds.Min.Y
}

func GetFloorCeilingAngle(room domain.RoomGeometry) float64 {
	floor, okFloor := geom.NormalizePlane(room.Floor)
	ceiling, okCeiling := geom.NormalizePlane(room.Ceiling)
	if !okFloor || !okCeiling {
		return 0
	}
	return geom.AngleBetween3(floor.Normal, ceiling.Normal)
}
