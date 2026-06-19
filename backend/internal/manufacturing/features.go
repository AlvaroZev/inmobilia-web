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

	finalFrontHeight := drawerHeight
	if !nestedInDesk && count == 1 {
		finalFrontHeight = volume.Height
	}

	subtractCarcass := !nestedInDesk && count == 1
	if feature.Params != nil {
		subtractCarcass = boolFromParams(feature.Params, "carcassFloorCeiling", subtractCarcass)
	}

	sideHeight := drawerBodySideHeight(finalFrontHeight, c.material.Thickness, subtractCarcass)

	boxDepth := drawerBoxDepth(volume.Depth, c.back)
	sideDepth := boxDepth
	runnerLength := snapRunnerLengthMm(sideDepth)

	boxW := drawerBoxWidth(volume.Width, nestedInDesk, c.material.Thickness)
	bottomWidth := c.drawerBottomWidth(boxW-2*c.material.Thickness, nestedInDesk)
	if nestedInDesk {
		bottomWidth = c.drawerBottomWidth(boxW-c.material.Thickness, nestedInDesk)
	}
	if bottomWidth < 0 {
		bottomWidth = c.drawerBottomWidth(c.innerWidth(volume)-2*c.material.Thickness, nestedInDesk)
	}
	bottomMaterialID := resolveBackMaterialIDFromParams(feature.Params)
	bottomMaterial := resolveMaterial(bottomMaterialID)
	bottomMaterial.Thickness = floatFromParams(feature.Params, "bottomThicknessMm", NordexThicknessMm)

	frontPanelWidth := drawerFrontPanelWidth(volume.Width)
	frontPanelHeight := drawerFrontPanelHeight(finalFrontHeight)
	if nestedInDesk {
		frontPanelWidth = drawerFrontPanelWidth(volume.Width - c.material.Thickness)
	}

	for i := 0; i < count; i++ {
		prefix := fmt.Sprintf("Cajón %s %d", volume.Label, i+1)
		var left, right domain.Part
		if !nestedInDesk {
			left = c.addPart(volume.ID, prefix+" lateral", string(domain.PartDrawerSide), sideDepth, sideHeight, c.material, "vertical")
			right = c.addPart(volume.ID, prefix+" lateral", string(domain.PartDrawerSide), sideDepth, sideHeight, c.material, "vertical")
		}
		bottom := c.addPart(volume.ID, prefix+" fondo", string(domain.PartDrawerBottom), bottomWidth, c.drawerBottomDepth(sideDepth), bottomMaterial, "horizontal")
		front := c.addPart(volume.ID, prefix+" frente", string(domain.PartDoor), frontPanelWidth, frontPanelHeight, c.material, "vertical")
		addPartThickTopEdgeBanding(c, front.ID, front.Width)
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
			Params:   json.RawMessage(fmt.Sprintf(`{"runner":"%s","lengthMm":%.0f}`, runner, runnerLength)),
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
