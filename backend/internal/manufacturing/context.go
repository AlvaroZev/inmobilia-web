package manufacturing

import (
	"fmt"

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
	return &compileContext{
		furnitureID: resolved.ID,
		root:        resolved.Root,
		material:    material,
		edgeBanding: resolveEdgeBanding(resolved.Root.MaterialID),
		back:        backMaterial(material),
	}
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
		Width:          width,
		Height:         height,
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
