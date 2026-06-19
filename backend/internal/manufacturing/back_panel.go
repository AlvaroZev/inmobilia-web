package manufacturing

import (
	"encoding/json"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

const (
	NordexThicknessMm         = 3
	MelamineBackThicknessMm   = 18
	PanelGrooveOffsetFromEdge = 18
	PanelGrooveWidthMm        = 4
	PanelGrooveDepthMm        = 7
	NordexGrooveInsetMm       = 6
	// Grooves are routed along the full panel edge (shop constraint).
)

func nordexPanelSpanMm(internalSpan float64, engagingGrooveCount int) float64 {
	return roundMm(internalSpan + float64(engagingGrooveCount)*NordexGrooveInsetMm)
}

func isNordexMaterialID(materialID string) bool {
	return materialID == "nordex"
}

func resolveBackMaterialIDFromParams(params json.RawMessage) string {
	if id := stringFromParams(params, "backMaterialId", ""); id != "" {
		return id
	}
	if id := stringFromParams(params, "bottomMaterialId", ""); id != "" {
		return id
	}
	return "melamine-white-18"
}

func resolveBackPanelMaterial(materialID string, board domain.Material) domain.Material {
	if isNordexMaterialID(materialID) {
		m := resolveMaterial("nordex")
		m.Thickness = NordexThicknessMm
		return m
	}
	m := board
	if m.Thickness < 10 {
		m.Thickness = MelamineBackThicknessMm
	}
	m.Name = m.Name + " (trasera)"
	m.Type = "back_panel"
	return m
}

func backPanelsUseGrooves(materialID string) bool {
	return isNordexMaterialID(materialID)
}

func carcassBackThicknessMm(materialID string, board domain.Material) float64 {
	return resolveBackPanelMaterial(materialID, board).Thickness
}
