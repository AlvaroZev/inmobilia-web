package solver

import (
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/roomgeometry"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
)

func computeInstallSpace(
	room domain.RoomGeometry,
	furniture domain.FurnitureDefinition,
	installation domain.InstallationConstraints,
) (box, error) {
	base, err := resolveBaseBounds(room, installation)
	if err != nil {
		return box{}, err
	}

	x := base.Min.X + installation.Clearances.Left
	y := base.Min.Y + installation.Clearances.Bottom
	z := base.Min.Z + installation.Clearances.Back

	floorOffset := installation.References.FloorOffset
	if furniture.Root.Adaptation != nil && furniture.Root.Adaptation.CompensateSkirting {
		skirting := detectSkirtingHeight(room, base)
		floorOffset = math.Max(floorOffset, skirting)
	}
	y += floorOffset

	width := (base.Max.X - base.Min.X) - installation.Clearances.Left - installation.Clearances.Right
	height := (base.Max.Y - base.Min.Y) - floorOffset - installation.Clearances.Top - installation.Clearances.Bottom
	depth := (base.Max.Z - base.Min.Z) - installation.Clearances.Back - installation.Clearances.Front

	if installation.References.CeilingOffset > 0 {
		height -= installation.References.CeilingOffset
	}

	width -= installation.Tolerances.Width
	height -= installation.Tolerances.Height
	depth -= installation.Tolerances.Depth

	if width <= geom.Epsilon || height <= geom.Epsilon || depth <= geom.Epsilon {
		return box{}, ErrInstallSpaceTooSmall
	}

	return box{X: x, Y: y, Z: z, Width: width, Height: height, Depth: depth}, nil
}

func resolveBaseBounds(room domain.RoomGeometry, installation domain.InstallationConstraints) (domain.BoundingBox3D, error) {
	if installation.Zone.Bounds != nil {
		return geom.NormalizeBounds(*installation.Zone.Bounds), nil
	}

	if installation.References.ReferenceWallID != "" {
		wall, ok := roomgeometry.FindWallByID(room, installation.References.ReferenceWallID)
		if !ok {
			return domain.BoundingBox3D{}, ErrReferenceWallNotFound
		}
		return wallToBounds(*wall, installation.References.WallOffset), nil
	}

	return roomPerimeterBounds(room)
}

func wallToBounds(wall domain.Wall, wallOffset float64) domain.BoundingBox3D {
	minX, minY, minZ := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	maxX, maxY, maxZ := -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

	for _, v := range wall.Vertices {
		minX = math.Min(minX, v.X)
		minY = math.Min(minY, v.Y)
		minZ = math.Min(minZ, v.Z)
		maxX = math.Max(maxX, v.X)
		maxY = math.Max(maxY, v.Y)
		maxZ = math.Max(maxZ, v.Z)
	}

	depth := wall.Thickness + wallOffset
	return domain.BoundingBox3D{
		Min: domain.Point3D{X: minX, Y: minY, Z: minZ},
		Max: domain.Point3D{X: maxX, Y: maxY, Z: minZ + depth},
	}
}

func roomPerimeterBounds(room domain.RoomGeometry) (domain.BoundingBox3D, error) {
	vertices := room.Perimeter.Vertices
	if len(vertices) < 3 {
		return domain.BoundingBox3D{}, ErrInstallZoneUndefined
	}

	minX, minZ := math.MaxFloat64, math.MaxFloat64
	maxX, maxZ := -math.MaxFloat64, -math.MaxFloat64

	for _, v := range vertices {
		minX = math.Min(minX, v.X)
		minZ = math.Min(minZ, v.Y)
		maxX = math.Max(maxX, v.X)
		maxZ = math.Max(maxZ, v.Y)
	}

	floorY := room.Floor.Point.Y
	ceilingY := room.Ceiling.Point.Y

	return domain.BoundingBox3D{
		Min: domain.Point3D{X: minX, Y: floorY, Z: minZ},
		Max: domain.Point3D{X: maxX, Y: ceilingY, Z: maxZ},
	}, nil
}

func detectSkirtingHeight(room domain.RoomGeometry, zone domain.BoundingBox3D) float64 {
	maxHeight := 0.0
	for _, obstacle := range roomgeometry.GetSkirtingObstacles(room) {
		if !boundsOverlapXZ(obstacle.Bounds, zone) {
			continue
		}
		h := roomgeometry.GetSkirtingHeight(obstacle)
		if h > maxHeight {
			maxHeight = h
		}
	}
	return maxHeight
}

func boundsOverlapXZ(a, b domain.BoundingBox3D) bool {
	a = geom.NormalizeBounds(a)
	b = geom.NormalizeBounds(b)
	return a.Min.X <= b.Max.X && a.Max.X >= b.Min.X &&
		a.Min.Z <= b.Max.Z && a.Max.Z >= b.Min.Z
}
