package geom

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

func PlaneFromPoints(a, b, c domain.Point3D) (domain.Plane, bool) {
	ab := Subtract3(b, a)
	ac := Subtract3(c, a)
	normal, ok := Normalize3(Cross3(ab, ac))
	if !ok {
		return domain.Plane{}, false
	}
	return domain.Plane{Point: a, Normal: normal}, true
}

func SignedDistanceToPlane(point domain.Point3D, plane domain.Plane) float64 {
	return Dot3(Subtract3(point, plane.Point), plane.Normal)
}

func ProjectPointOnPlane(point domain.Point3D, plane domain.Plane) domain.Point3D {
	dist := SignedDistanceToPlane(point, plane)
	return Subtract3(point, Scale3(plane.Normal, dist))
}

func NormalizePlane(plane domain.Plane) (domain.Plane, bool) {
	normal, ok := Normalize3(plane.Normal)
	if !ok {
		return domain.Plane{}, false
	}
	return domain.Plane{Point: plane.Point, Normal: normal}, true
}
