package solver

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

// box is an axis-aligned installation or volume region in millimeters.
type box struct {
	X      float64
	Y      float64
	Z      float64
	Width  float64
	Height float64
	Depth  float64
}

func (b box) dimension(axis domain.SplitAxis) float64 {
	switch axis {
	case domain.SplitAxisX:
		return b.Width
	case domain.SplitAxisY:
		return b.Height
	case domain.SplitAxisZ:
		return b.Depth
	default:
		return 0
	}
}

func (b box) setDimension(axis domain.SplitAxis, value float64) box {
	switch axis {
	case domain.SplitAxisX:
		b.Width = value
	case domain.SplitAxisY:
		b.Height = value
	case domain.SplitAxisZ:
		b.Depth = value
	}
	return b
}

func (b box) setOriginAxis(axis domain.SplitAxis, value float64) box {
	switch axis {
	case domain.SplitAxisX:
		b.X = value
	case domain.SplitAxisY:
		b.Y = value
	case domain.SplitAxisZ:
		b.Z = value
	}
	return b
}

func (b box) toBounds() domain.BoundingBox3D {
	return domain.BoundingBox3D{
		Min: domain.Point3D{X: b.X, Y: b.Y, Z: b.Z},
		Max: domain.Point3D{
			X: b.X + b.Width,
			Y: b.Y + b.Height,
			Z: b.Z + b.Depth,
		},
	}
}
