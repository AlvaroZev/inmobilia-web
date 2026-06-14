package manufacturing

import (
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/resolvedfurniture"
)

func isNestedDrawerTower(ref resolvedfurniture.VolumeRef) bool {
	if ref.Parent == nil {
		return false
	}
	hasDrawer := false
	for _, feature := range ref.Volume.Features {
		if feature.Type == "drawer_stack" {
			hasDrawer = true
			break
		}
	}
	return hasDrawer && volumeHasDeskFrame(*ref.Parent)
}

func volumeHasDeskFrame(volume domain.ResolvedVolume) bool {
	for _, feature := range volume.Features {
		if feature.Type == "desk_frame" {
			return true
		}
	}
	return false
}

func isDeskLayout(volume domain.ResolvedVolume) bool {
	if volumeHasDeskFrame(volume) {
		return true
	}
	for _, child := range volume.Children {
		if isDeskLayout(child) {
			return true
		}
	}
	return false
}

// CompileManufacturing transforms resolved furniture into fabrication parts.
func CompileManufacturing(resolved domain.ResolvedFurniture) (domain.ManufacturingModel, error) {
	if result := resolvedfurniture.ValidateResolvedFurniture(resolved); !result.Valid {
		return domain.ManufacturingModel{}, ErrInvalidResolvedFurniture
	}

	ctx := newCompileContext(resolved)

	resolvedfurniture.WalkResolvedTree(resolved.Root, func(ref resolvedfurniture.VolumeRef) bool {
		if ref.Parent == nil && !isDeskLayout(ref.Volume) {
			compileOuterCarcass(ctx, ref.Volume)
		}

		if len(ref.Volume.Children) > 0 {
			compileDividers(ctx, ref.Volume)
		}

		if isNestedDrawerTower(ref) {
			compileNestedDrawerTower(ctx, ref.Volume)
		}

		compileFeatures(ctx, ref.Volume)

		compileFronts(ctx, ref.Volume)
		return true
	})

	model := ctx.model()
	if result := ValidateManufacturingModel(model); !result.Valid {
		return domain.ManufacturingModel{}, ErrInvalidManufacturingModel
	}

	return model, nil
}
