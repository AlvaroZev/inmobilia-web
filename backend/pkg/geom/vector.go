package geom

import (
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

const Epsilon = 0.001

func Subtract3(a, b domain.Point3D) domain.Point3D {
	return domain.Point3D{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func Add3(a, b domain.Point3D) domain.Point3D {
	return domain.Point3D{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

func Scale3(v domain.Point3D, s float64) domain.Point3D {
	return domain.Point3D{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

func Dot3(a, b domain.Point3D) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Cross3(a, b domain.Point3D) domain.Point3D {
	return domain.Point3D{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func Magnitude3(v domain.Point3D) float64 {
	return math.Hypot(v.X, math.Hypot(v.Y, v.Z))
}

func Normalize3(v domain.Point3D) (domain.Point3D, bool) {
	len := Magnitude3(v)
	if len < Epsilon {
		return domain.Point3D{}, false
	}
	return Scale3(v, 1/len), true
}

func Distance3(a, b domain.Point3D) float64 {
	return Magnitude3(Subtract3(a, b))
}

func Distance2(a, b domain.Point2D) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func AngleBetween3(a, b domain.Point3D) float64 {
	na, okA := Normalize3(a)
	nb, okB := Normalize3(b)
	if !okA || !okB {
		return 0
	}
	cos := Dot3(na, nb)
	cos = math.Min(1, math.Max(-1, cos))
	return math.Acos(cos)
}

func AngleBetween2(a, b domain.Point2D) float64 {
	lenA := math.Hypot(a.X, a.Y)
	lenB := math.Hypot(b.X, b.Y)
	if lenA < Epsilon || lenB < Epsilon {
		return 0
	}
	cos := (a.X*b.X + a.Y*b.Y) / (lenA * lenB)
	cos = math.Min(1, math.Max(-1, cos))
	return math.Acos(cos)
}
