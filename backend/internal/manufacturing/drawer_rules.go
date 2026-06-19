package manufacturing

import (
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

// Drawer / front manufacturing rules (mirrors frontend drawer-config).
const (
	ThickEdgeBandingMm       = 3
	FrontToleranceInsetMm      = 1
	CarcassBackThicknessMm     = 18
	DrawerRearClearanceMm      = 10
	LateralInsetMm             = 13.5
	DrawerWidthReductionMm     = 27
	BodySideHeightRatio        = 0.85
	FalseFrontHeightRatio      = 0.85
	BodyTopClearanceMm         = 24
	RunnerLengthMinMm          = 250
	RunnerLengthMaxMm          = 650
	RunnerLengthStepMm         = 50
	RunnerOffsetFromBaseMm     = 50
)

func drawerInteriorAvailableHeight(slotHeight float64, panelThickness float64, subtractCarcass bool) float64 {
	inset := 0.0
	if subtractCarcass {
		inset = 2 * panelThickness
	}
	return math.Max(0, slotHeight-inset)
}

func drawerInteriorMaxBodyHeight(slotHeight float64, panelThickness float64) float64 {
	inner := slotHeight - 2*panelThickness - BodyTopClearanceMm
	return math.Max(0, inner)
}

func drawerBodySideHeight(slotHeight float64, panelThickness float64, subtractCarcass bool) float64 {
	available := drawerInteriorAvailableHeight(slotHeight, panelThickness, subtractCarcass)
	h := available * BodySideHeightRatio
	if subtractCarcass {
		h = math.Min(h, drawerInteriorMaxBodyHeight(slotHeight, panelThickness))
	}
	return roundMm(math.Max(0, h))
}

func drawerFalseFrontHeight(bodySideHeight float64) float64 {
	return roundMm(math.Max(0, bodySideHeight*FalseFrontHeightRatio))
}

func drawerFrontPanelWidth(externalWidth float64) float64 {
	return externalWidth - 2*FrontToleranceInsetMm - 2*ThickEdgeBandingMm
}

func drawerFrontPanelHeight(externalHeight float64) float64 {
	return externalHeight - 2*FrontToleranceInsetMm - 2*ThickEdgeBandingMm
}

func carcassStructureDepth(volumeDepth float64, back domain.Material) float64 {
	if backPanelsUseGrooves(back.ID) {
		return volumeDepth
	}
	return volumeDepth - back.Thickness
}

func drawerBoxDepth(volumeDepth float64, back domain.Material) float64 {
	return carcassStructureDepth(volumeDepth, back) - DrawerRearClearanceMm
}

func drawerBoxWidth(externalWidth float64, nestedInDesk bool, panelThickness float64) float64 {
	innerW := externalWidth - 2*panelThickness
	if nestedInDesk {
		return innerW - LateralInsetMm
	}
	return innerW - DrawerWidthReductionMm
}

func snapRunnerLengthMm(availableDepth float64) float64 {
	best := float64(RunnerLengthMinMm)
	for length := RunnerLengthMinMm; length <= RunnerLengthMaxMm; length += RunnerLengthStepMm {
		if float64(length) <= availableDepth {
			best = float64(length)
			continue
		}
		break
	}
	return best
}
