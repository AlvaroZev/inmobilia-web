package resolvedfurniture

import (
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
)

func VolumeBounds(volume domain.ResolvedVolume) domain.BoundingBox3D {
	return domain.BoundingBox3D{
		Min: domain.Point3D{X: volume.X, Y: volume.Y, Z: volume.Z},
		Max: domain.Point3D{
			X: volume.X + volume.Width,
			Y: volume.Y + volume.Height,
			Z: volume.Z + volume.Depth,
		},
	}
}

func FurnitureBounds(furniture domain.ResolvedFurniture) domain.BoundingBox3D {
	return VolumeBounds(furniture.Root)
}

func VolumeContains(outer, inner domain.ResolvedVolume) bool {
	ob := VolumeBounds(outer)
	ib := VolumeBounds(inner)

	return ib.Min.X >= ob.Min.X-geom.Epsilon &&
		ib.Min.Y >= ob.Min.Y-geom.Epsilon &&
		ib.Min.Z >= ob.Min.Z-geom.Epsilon &&
		ib.Max.X <= ob.Max.X+geom.Epsilon &&
		ib.Max.Y <= ob.Max.Y+geom.Epsilon &&
		ib.Max.Z <= ob.Max.Z+geom.Epsilon
}

func BoundsOverlap(a, b domain.BoundingBox3D) bool {
	a = geom.NormalizeBounds(a)
	b = geom.NormalizeBounds(b)

	return a.Min.X < b.Max.X-geom.Epsilon &&
		a.Max.X > b.Min.X+geom.Epsilon &&
		a.Min.Y < b.Max.Y-geom.Epsilon &&
		a.Max.Y > b.Min.Y+geom.Epsilon &&
		a.Min.Z < b.Max.Z-geom.Epsilon &&
		a.Max.Z > b.Min.Z+geom.Epsilon
}

func VolumeOverlap(a, b domain.ResolvedVolume) bool {
	return BoundsOverlap(VolumeBounds(a), VolumeBounds(b))
}

func CompartmentVolumeMm3(volume domain.ResolvedVolume) float64 {
	return volume.Width * volume.Height * volume.Depth
}

func TotalLeafVolumeMm3(root domain.ResolvedVolume) float64 {
	var total float64
	for _, leaf := range GetLeafVolumes(root) {
		total += CompartmentVolumeMm3(leaf)
	}
	return total
}

func ExternalDimensions(furniture domain.ResolvedFurniture) (width, height, depth float64) {
	root := furniture.Root
	return root.Width, root.Height, root.Depth
}
