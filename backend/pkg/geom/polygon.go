package geom

import (
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func PolygonArea2D(polygon domain.Polygon2D) float64 {
	vertices := polygon.Vertices
	if len(vertices) < 3 {
		return 0
	}

	var sum float64
	for i := range vertices {
		j := (i + 1) % len(vertices)
		sum += vertices[i].X*vertices[j].Y - vertices[j].X*vertices[i].Y
	}
	return sum / 2
}

func PolygonPerimeter2D(polygon domain.Polygon2D) float64 {
	vertices := polygon.Vertices
	if len(vertices) < 2 {
		return 0
	}

	var perimeter float64
	for i := range vertices {
		j := (i + 1) % len(vertices)
		perimeter += Distance2(vertices[i], vertices[j])
	}
	return perimeter
}

func InteriorAngleAtVertex2D(polygon domain.Polygon2D, index int) float64 {
	vertices := polygon.Vertices
	n := len(vertices)
	if n < 3 || index < 0 || index >= n {
		return 0
	}

	prev := vertices[(index-1+n)%n]
	curr := vertices[index]
	next := vertices[(index+1)%n]

	incoming := domain.Point2D{X: prev.X - curr.X, Y: prev.Y - curr.Y}
	outgoing := domain.Point2D{X: next.X - curr.X, Y: next.Y - curr.Y}

	return math.Pi - AngleBetween2(incoming, outgoing)
}

func NormalizeBounds(bounds domain.BoundingBox3D) domain.BoundingBox3D {
	return domain.BoundingBox3D{
		Min: domain.Point3D{
			X: math.Min(bounds.Min.X, bounds.Max.X),
			Y: math.Min(bounds.Min.Y, bounds.Max.Y),
			Z: math.Min(bounds.Min.Z, bounds.Max.Z),
		},
		Max: domain.Point3D{
			X: math.Max(bounds.Min.X, bounds.Max.X),
			Y: math.Max(bounds.Min.Y, bounds.Max.Y),
			Z: math.Max(bounds.Min.Z, bounds.Max.Z),
		},
	}
}

func BoundsVolume(bounds domain.BoundingBox3D) float64 {
	b := NormalizeBounds(bounds)
	return (b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y) * (b.Max.Z - b.Min.Z)
}
