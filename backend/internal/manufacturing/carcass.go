package manufacturing

import (
	"fmt"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func compileOuterCarcass(c *compileContext, volume domain.ResolvedVolume) {
	t := c.material.Thickness
	innerW := c.innerWidth(volume)
	innerH := c.innerHeight(volume)
	structureDepth := carcassStructureDepth(volume.Depth, c.back)

	left := c.addPart(volume.ID, "Lateral izquierdo", string(domain.PartLateral), structureDepth, volume.Height, c.material, "vertical")
	right := c.addPart(volume.ID, "Lateral derecho", string(domain.PartLateral), structureDepth, volume.Height, c.material, "vertical")
	base := c.addPart(volume.ID, "Base", string(domain.PartBase), innerW, structureDepth, c.material, "horizontal")
	top := c.addPart(volume.ID, "Techo", string(domain.PartTop), innerW, structureDepth, c.material, "horizontal")
	back := c.addPart(volume.ID, "Trasera", string(domain.PartBack), c.backPanelWidth(innerW), c.backPanelHeight(innerH), c.back, "")

	addPartEdgeBanding(c, left.ID, domain.EdgeTop, left.Width)
	addPartEdgeBanding(c, left.ID, domain.EdgeBottom, left.Width)
	addPartEdgeBanding(c, right.ID, domain.EdgeTop, right.Width)
	addPartEdgeBanding(c, right.ID, domain.EdgeBottom, right.Width)
	addPartEdgeBanding(c, base.ID, domain.EdgeTop, base.Width)
	addPartEdgeBanding(c, top.ID, domain.EdgeTop, top.Width)

	_ = t
	_ = back
}

func compileNestedDrawerTower(c *compileContext, volume domain.ResolvedVolume) {
	if !nestedDrawerTowerHasBase(volume) {
		return
	}
	innerW := volume.Width - 2*c.material.Thickness
	base := c.addPart(volume.ID, "Base cajonera", string(domain.PartBase), innerW, volume.Depth, c.material, "horizontal")
	addPartEdgeBanding(c, base.ID, domain.EdgeTop, base.Width)
}

func nestedDrawerTowerHasBase(volume domain.ResolvedVolume) bool {
	for _, feature := range volume.Features {
		if feature.Type != "drawer_stack" {
			continue
		}
		return boolFromParams(feature.Params, "hasBase", true)
	}
	return true
}

func compileDividers(c *compileContext, parent domain.ResolvedVolume) {
	if len(parent.Children) < 2 {
		return
	}

	for i := 0; i < len(parent.Children)-1; i++ {
		left := parent.Children[i]
		right := parent.Children[i+1]

		switch axisFromChildren(left, right) {
		case domain.SplitAxisX:
			divider := c.addPart(
				parent.ID,
				fmt.Sprintf("División vertical %d", i+1),
				string(domain.PartDivider),
				parent.Depth,
				parent.Height,
				c.material,
				"vertical",
			)
			addPartEdgeBanding(c, divider.ID, domain.EdgeTop, divider.Width)
		case domain.SplitAxisY:
			divider := c.addPart(
				parent.ID,
				fmt.Sprintf("División horizontal %d", i+1),
				string(domain.PartDivider),
				c.innerWidth(parent),
				parent.Depth,
				c.material,
				"horizontal",
			)
			addPartEdgeBanding(c, divider.ID, domain.EdgeTop, divider.Width)
		case domain.SplitAxisZ:
			divider := c.addPart(
				parent.ID,
				fmt.Sprintf("División profundidad %d", i+1),
				string(domain.PartDivider),
				c.innerWidth(parent),
				parent.Height,
				c.material,
				"vertical",
			)
			addPartEdgeBanding(c, divider.ID, domain.EdgeRight, divider.Height)
		}
	}
}

func axisFromChildren(left, right domain.ResolvedVolume) domain.SplitAxis {
	if left.X != right.X {
		return domain.SplitAxisX
	}
	if left.Y != right.Y {
		return domain.SplitAxisY
	}
	return domain.SplitAxisZ
}
