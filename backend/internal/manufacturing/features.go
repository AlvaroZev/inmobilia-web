package manufacturing

import (
	"encoding/json"
	"fmt"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func compileFeatures(c *compileContext, volume domain.ResolvedVolume) {
	for _, feature := range volume.Features {
		switch feature.Type {
		case "desk_frame":
			compileDeskFrame(c, volume, feature)
		case "shelf_set":
			compileShelfSet(c, volume, feature)
		case "drawer_stack":
			compileDrawerStack(c, volume, feature)
		case "hanger_rod":
			compileHangerRod(c, volume, feature)
		}
	}
}

func compileDeskFrame(c *compileContext, volume domain.ResolvedVolume, feature domain.ResolvedFeature) {
	overhang := floatFromParams(feature.Params, "topOverhangMm", 25)
	braceRatio := floatFromParams(feature.Params, "braceHeightRatio", 0.5)
	if braceRatio <= 0 || braceRatio > 1 {
		braceRatio = 0.5
	}

	legHeight := volume.Height
	braceHeight := legHeight * braceRatio
	topWidth := volume.Width + 2*overhang
	topDepth := volume.Depth + overhang

	left := c.addPart(volume.ID, "Lateral izquierdo", string(domain.PartLateral), volume.Depth, legHeight, c.material, "vertical")
	right := c.addPart(volume.ID, "Lateral derecho", string(domain.PartLateral), volume.Depth, legHeight, c.material, "vertical")
	brace := c.addPart(volume.ID, "Amarre trasero", string(domain.PartBack), volume.Width-2*c.material.Thickness, braceHeight, c.material, "vertical")
	top := c.addPart(volume.ID, "Mesa escritorio", string(domain.PartTop), topWidth, topDepth, c.material, "horizontal")

	addPartEdgeBanding(c, left.ID, domain.EdgeTop, left.Width)
	addPartEdgeBanding(c, left.ID, domain.EdgeBottom, left.Width)
	addPartEdgeBanding(c, right.ID, domain.EdgeTop, right.Width)
	addPartEdgeBanding(c, right.ID, domain.EdgeBottom, right.Width)
	addPartEdgeBanding(c, brace.ID, domain.EdgeTop, brace.Width)
	addPartEdgeBanding(c, top.ID, domain.EdgeTop, top.Width)
}

func compileShelfSet(c *compileContext, volume domain.ResolvedVolume, feature domain.ResolvedFeature) {
	count := intFromParams(feature.Params, "count", 1)
	shelfWidth := c.innerWidth(volume)
	shelfDepth := c.innerDepth(volume)

	for i := 0; i < count; i++ {
		part := c.addPart(
			volume.ID,
			fmt.Sprintf("Repisa %s %d", volume.Label, i+1),
			string(domain.PartShelf),
			shelfWidth,
			shelfDepth,
			c.material,
			"horizontal",
		)
		addPartEdgeBanding(c, part.ID, domain.EdgeTop, part.Width)
	}
}

func compileDrawerStack(c *compileContext, volume domain.ResolvedVolume, feature domain.ResolvedFeature) {
	count := intFromParams(feature.Params, "count", 1)
	runner := stringFromParams(feature.Params, "runner", "standard")
	sharedLateral := stringFromParams(feature.Params, "sharedLateral", "")
	nestedInDesk := sharedLateral == "right"

	drawerHeight := floatFromParams(feature.Params, "drawerHeightMm", 0)
	if nestedInDesk {
		if drawerHeight <= 0 {
			drawerHeight = 175
		}
	} else if drawerHeight <= 0 {
		drawerHeight = volume.Height / float64(count)
	}

	sideHeight := drawerHeight - 4
	if sideHeight < 40 {
		sideHeight = drawerHeight
	}
	sideDepth := c.innerDepth(volume)
	bottomWidth := c.innerWidth(volume) - 2*c.material.Thickness
	if nestedInDesk {
		bottomWidth = volume.Width - 2*c.material.Thickness
	}
	bottomMaterial := resolveMaterial(stringFromParams(feature.Params, "bottomMaterialId", "nordex"))
	bottomMaterial.Thickness = floatFromParams(feature.Params, "bottomThicknessMm", 3)
	if !nestedInDesk {
		bottomMaterial = c.material
	}

	frontWidth := volume.Width
	if nestedInDesk {
		frontWidth = volume.Width - c.material.Thickness
	}
	frontHeight := drawerHeight

	for i := 0; i < count; i++ {
		prefix := fmt.Sprintf("Cajón %s %d", volume.Label, i+1)
		var left, right domain.Part
		if !nestedInDesk {
			left = c.addPart(volume.ID, prefix+" lateral", string(domain.PartDrawerSide), sideDepth, sideHeight, c.material, "vertical")
			right = c.addPart(volume.ID, prefix+" lateral", string(domain.PartDrawerSide), sideDepth, sideHeight, c.material, "vertical")
		}
		bottom := c.addPart(volume.ID, prefix+" fondo", string(domain.PartDrawerBottom), bottomWidth, sideDepth, bottomMaterial, "horizontal")
		front := c.addPart(volume.ID, prefix+" frente", string(domain.PartDoor), frontWidth, frontHeight-2, c.material, "vertical")
		addPartEdgeBanding(c, front.ID, domain.EdgeTop, front.Width)
		addPartEdgeBanding(c, front.ID, domain.EdgeBottom, front.Width)
		addPartEdgeBanding(c, front.ID, domain.EdgeLeft, front.Height)
		addPartEdgeBanding(c, front.ID, domain.EdgeRight, front.Height)

		if !nestedInDesk {
			addPartEdgeBanding(c, left.ID, domain.EdgeTop, left.Width)
			addPartEdgeBanding(c, right.ID, domain.EdgeTop, right.Width)
		}

		partIDs := []string{bottom.ID, front.ID}
		if !nestedInDesk {
			partIDs = []string{left.ID, right.ID, bottom.ID}
		}

		c.hardware = append(c.hardware, domain.Hardware{
			ID:       volume.ID + "-runner-" + fmt.Sprint(i+1),
			Type:     "drawer_runner",
			Quantity: 1,
			PartIDs:  partIDs,
			Params:   json.RawMessage(`{"runner":"` + runner + `"}`),
		})
	}
}

func compileHangerRod(c *compileContext, volume domain.ResolvedVolume, feature domain.ResolvedFeature) {
	c.hardware = append(c.hardware,
		domain.Hardware{
			ID:       feature.ID + "-rod",
			Type:     "hanger_rod",
			Quantity: 1,
			Params:   json.RawMessage(`{"length":` + fmt.Sprint(c.innerWidth(volume)) + `}`),
		},
		domain.Hardware{
			ID:       feature.ID + "-bracket",
			Type:     "rod_bracket",
			Quantity: 2,
			PartIDs:  []string{feature.ID},
		},
	)
}

func intFromParams(params json.RawMessage, key string, fallback int) int {
	var data map[string]any
	if err := json.Unmarshal(params, &data); err != nil {
		return fallback
	}
	value, ok := data[key]
	if !ok {
		return fallback
	}
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return fallback
	}
}

func boolFromParams(params json.RawMessage, key string, fallback bool) bool {
	var data map[string]any
	if err := json.Unmarshal(params, &data); err != nil {
		return fallback
	}
	value, ok := data[key]
	if !ok {
		return fallback
	}
	switch v := value.(type) {
	case bool:
		return v
	default:
		return fallback
	}
}

func floatFromParams(params json.RawMessage, key string, fallback float64) float64 {
	var data map[string]any
	if err := json.Unmarshal(params, &data); err != nil {
		return fallback
	}
	value, ok := data[key]
	if !ok {
		return fallback
	}
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		return fallback
	}
}

func stringFromParams(params json.RawMessage, key, fallback string) string {
	var data map[string]any
	if err := json.Unmarshal(params, &data); err != nil {
		return fallback
	}
	value, ok := data[key].(string)
	if !ok {
		return fallback
	}
	return value
}
