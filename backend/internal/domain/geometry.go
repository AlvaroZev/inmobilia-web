package domain

// All dimensions in millimeters. Angles are never stored — derive from geometry.

type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Plane is defined by a point on the surface and its outward normal.
type Plane struct {
	Point  Point3D `json:"point"`
	Normal Point3D `json:"normal"`
}

type Polygon2D struct {
	Vertices []Point2D `json:"vertices"`
}

type BoundingBox3D struct {
	Min Point3D `json:"min"`
	Max Point3D `json:"max"`
}
