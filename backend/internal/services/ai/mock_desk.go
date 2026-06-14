package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

const deskFrameParams = `{"braceHeightRatio":0.5,"topOverhangMm":25}`

const drawerStackBaseParams = `"sharedLateral":"right","drawerHeightMm":175,"backClearanceMm":40,"bottomMaterialId":"nordex","bottomThicknessMm":3,"grooveWidthMm":18,"grooveDepthMm":7,"grooveRailThicknessMm":4,"runnerHeightMm":40,"runnerWidthMm":8,"runnerLengthStepMm":50,"runnerLengthMinMm":200,"boxInsetSideMm":2`

func deskDrawerStackParams(intent deskDrawerIntent) json.RawMessage {
	hasBase := intent.mode == "tower"
	raw := fmt.Sprintf(
		`{"count":%d,"runner":"soft-close","drawerMode":"%s",%s,"hasBase":%t}`,
		intent.count,
		intent.mode,
		drawerStackBaseParams,
		hasBase,
	)
	return json.RawMessage(raw)
}

func (p *MockParser) deskDefinition(description, name string) domain.FurnitureDefinition {
	if name == "" {
		name = "Escritorio"
	}

	lower := strings.ToLower(description)
	drawerIntent := parseDeskDrawerIntent(description)
	depth := 600.0
	if strings.Contains(lower, "450") {
		depth = 450
	}

	deskFrame := domain.Feature{
		ID:     "desk-structure",
		Type:   "desk_frame",
		Params: json.RawMessage(deskFrameParams),
	}

	root := domain.VolumeNode{
		ID:    "root",
		Label: "Escritorio",
		Constraints: domain.VolumeConstraints{
			Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
			Height: &domain.DimensionConstraint{Mode: domain.DimensionFixed, Value: 720},
			Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFixed, Value: depth},
		},
		Features: []domain.Feature{deskFrame},
		Fronts:   []domain.Front{},
		Adaptation: &domain.AdaptationRules{
			FollowFloor:   true,
			FollowCeiling: false,
		},
		Manufacturing: &domain.ManufacturingHints{
			MaterialID:  "melamine-white-18",
			EdgeBanding: "pvc-white-1mm",
			BackPanel:   false,
		},
	}

	if drawerIntent.enabled {
		bayID := "drawer-bay"
		bayLabel := "Cajón"
		if drawerIntent.mode == "tower" {
			bayID = "drawer-tower"
			bayLabel = "Torre de cajones"
		}

		root.Split = &domain.VolumeSplit{
			Axis:   domain.SplitAxisX,
			Ratios: []float64{drawerIntent.legRatio, drawerIntent.bayRatio},
		}
		root.Children = []domain.VolumeNode{
			{
				ID:    "leg-space",
				Label: "Espacio de piernas",
				Constraints: domain.VolumeConstraints{
					Width:  &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: drawerIntent.legRatio},
					Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
					Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				},
				Children: []domain.VolumeNode{},
				Features: []domain.Feature{},
				Fronts:   []domain.Front{},
			},
			{
				ID:    bayID,
				Label: bayLabel,
				Constraints: domain.VolumeConstraints{
					Width:  &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: drawerIntent.bayRatio},
					Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
					Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				},
				Children: []domain.VolumeNode{},
				Features: []domain.Feature{
					{
						ID:     "desk-drawers",
						Type:   "drawer_stack",
						Params: deskDrawerStackParams(drawerIntent),
					},
				},
				Fronts: []domain.Front{},
			},
		}
	} else {
		root.Children = []domain.VolumeNode{}
	}

	return domain.FurnitureDefinition{
		ID:          fmt.Sprintf("ai-desk-%d", time.Now().Unix()),
		Name:        name,
		Description: description,
		Root:        root,
	}
}

func defaultManufacturingHints() *domain.ManufacturingHints {
	return &domain.ManufacturingHints{
		MaterialID:  "melamine-white-18",
		EdgeBanding: "pvc-white-1mm",
		BackPanel:   true,
	}
}
