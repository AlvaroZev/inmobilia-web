package manufacturing

import (
	"fmt"
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

type compileContext struct {
	furnitureID string
	root        domain.ResolvedVolume
	material    domain.Material
	edgeBanding string
	back        domain.Material

	parts       []domain.Part
	hardware    []domain.Hardware
	edgeBandingList []domain.EdgeBanding
	drilling    []domain.Drilling

	partCounter int
}

func newCompileContext(resolved domain.ResolvedFurniture) *compileContext {
	material := resolveMaterial(resolved.Root.MaterialID)
	backID := resolveVolumeBackMaterialID(resolved.Root)
	return &compileContext{
		furnitureID: resolved.ID,
		root:        resolved.Root,
		material:    material,
		edgeBanding: resolveEdgeBanding(resolved.Root.MaterialID),
		back:        resolveBackPanelMaterial(backID, material),
	}
}

func resolveVolumeBackMaterialID(volume domain.ResolvedVolume) string {
	for _, feature := range volume.Features {
		if feature.Type == "drawer_stack" {
			return resolveBackMaterialIDFromParams(feature.Params)
		}
	}
	return "melamine-white-18"
}

func (c *compileContext) nextPartID(volumeID, partType string) string {
	c.partCounter++
	return fmt.Sprintf("%s-%s-%d", volumeID, partType, c.partCounter)
}

func (c *compileContext) addPart(volumeID, name, partType string, width, height float64, material domain.Material, grain string) domain.Part {
	part := domain.Part{
		ID:             c.nextPartID(volumeID, partType),
		Name:           name,
		Type:           partType,
		VolumeID:       volumeID,
		Width:          roundMm(width),
		Height:         roundMm(height),
		Thickness:      material.Thickness,
		Material:       material,
		GrainDirection: grain,
	}
	c.parts = append(c.parts, part)
	return part
}

func (c *compileContext) innerWidth(volume domain.ResolvedVolume) float64 {
	return volume.Width - 2*c.material.Thickness
}

func (c *compileContext) innerHeight(volume domain.ResolvedVolume) float64 {
	return volume.Height - 2*c.material.Thickness
}

func (c *compileContext) backPanelWidth(innerW float64) float64 {
	if backPanelsUseGrooves(c.back.ID) {
		return nordexPanelSpanMm(innerW, 2)
	}
	return innerW
}

func (c *compileContext) backPanelHeight(innerH float64) float64 {
	if backPanelsUseGrooves(c.back.ID) {
		return nordexPanelSpanMm(innerH, 2)
	}
	return innerH
}

func (c *compileContext) drawerBottomWidth(innerW float64, nestedInDesk bool) float64 {
	if backPanelsUseGrooves(resolveBackMaterialIDFromVolume(c.root)) {
		grooves := 2
		if nestedInDesk {
			grooves = 1
		}
		return nordexPanelSpanMm(innerW, grooves)
	}
	return innerW
}

func (c *compileContext) drawerBottomDepth(sideDepth float64) float64 {
	innerD := sideDepth - 2*c.material.Thickness
	if innerD < 0 {
		innerD = 0
	}
	if backPanelsUseGrooves(resolveBackMaterialIDFromVolume(c.root)) {
		return nordexPanelSpanMm(innerD, 2)
	}
	return math.Max(0, sideDepth-c.material.Thickness)
}

func resolveBackMaterialIDFromVolume(volume domain.ResolvedVolume) string {
	for _, feature := range volume.Features {
		if feature.Type == "drawer_stack" {
			return resolveBackMaterialIDFromParams(feature.Params)
		}
	}
	return "melamine-white-18"
}

func (c *compileContext) innerDepth(volume domain.ResolvedVolume) float64 {
	depth := volume.Depth - c.back.Thickness - 8
	if depth < c.material.Thickness {
		depth = volume.Depth - c.back.Thickness
	}
	return depth
}

func (c *compileContext) model() domain.ManufacturingModel {
	return domain.ManufacturingModel{
		FurnitureID: c.furnitureID,
		Parts:       c.parts,
		Hardware:    c.hardware,
		EdgeBanding: c.edgeBandingList,
		Drilling:    c.drilling,
	}
}
