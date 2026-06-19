package manufacturing

import (
	"encoding/json"
	"fmt"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func compileFronts(c *compileContext, volume domain.ResolvedVolume) {
	for _, front := range volume.Fronts {
		switch front.Type {
		case "door":
			compileDoor(c, volume, front)
		case "drawer_front":
			compileDrawerFronts(c, volume, front)
		}
	}
}

func compileDoor(c *compileContext, volume domain.ResolvedVolume, front domain.ResolvedFront) {
	part := c.addPart(
		volume.ID,
		fmt.Sprintf("Puerta %s", volume.Label),
		string(domain.PartDoor),
		front.Width,
		front.Height,
		c.material,
		"vertical",
	)

	addPartThickTopEdgeBanding(c, part.ID, part.Width)
	addPartEdgeBanding(c, part.ID, domain.EdgeBottom, part.Width)
	addPartEdgeBanding(c, part.ID, domain.EdgeLeft, part.Height)
	addPartEdgeBanding(c, part.ID, domain.EdgeRight, part.Height)

	hingeCount := 2
	if front.Height > 2000 {
		hingeCount = 3
	}

	hinge := stringFromParams(front.Params, "hinge", "standard")
	c.hardware = append(c.hardware, domain.Hardware{
		ID:       front.ID + "-hinges",
		Type:     "hinge",
		Quantity: hingeCount,
		PartIDs:  []string{part.ID},
		Params:   json.RawMessage(`{"hinge":"` + hinge + `"}`),
	})
}

func compileDrawerFronts(c *compileContext, volume domain.ResolvedVolume, front domain.ResolvedFront) {
	drawerCount := drawerCountFromVolume(volume)
	if drawerCount < 1 {
		drawerCount = 1
	}

	frontHeight := front.Height / float64(drawerCount)
	for i := 0; i < drawerCount; i++ {
		part := c.addPart(
			volume.ID,
			fmt.Sprintf("Frente cajón %s %d", volume.Label, i+1),
			string(domain.PartDoor),
			front.Width,
			frontHeight,
			c.material,
			"vertical",
		)
		addPartEdgeBanding(c, part.ID, domain.EdgeTop, part.Width)
		addPartEdgeBanding(c, part.ID, domain.EdgeBottom, part.Width)
		addPartEdgeBanding(c, part.ID, domain.EdgeLeft, part.Height)
		addPartEdgeBanding(c, part.ID, domain.EdgeRight, part.Height)
	}
}

func drawerCountFromVolume(volume domain.ResolvedVolume) int {
	for _, feature := range volume.Features {
		if feature.Type == "drawer_stack" {
			return intFromParams(feature.Params, "count", 1)
		}
	}
	return 1
}
