package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func (p *MockParser) entertainmentCenterDefinition(description, name string) domain.FurnitureDefinition {
	if name == "" {
		name = "Centro de entretenimiento"
	}

	lower := strings.ToLower(description)
	bodyCount := parseBodyCount(description)
	if bodyCount < 2 {
		bodyCount = 2
	}
	if bodyCount > 4 {
		bodyCount = 4
	}

	ratios := equalRatios(bodyCount)
	depth := 450.0
	if strings.Contains(lower, "600") {
		depth = 600
	}

	withGlass := strings.Contains(lower, "vidrio") || strings.Contains(lower, "glass")

	tvFront := domain.Front{}
	if withGlass {
		tvFront = domain.Front{
			ID:     "tv-glass",
			Type:   "glass",
			Params: json.RawMessage(`{"materialId":"clear-glass"}`),
		}
	}

	storageChildren := make([]domain.VolumeNode, bodyCount)
	for i := range bodyCount {
		id := fmt.Sprintf("storage-%d", i+1)
		storageChildren[i] = domain.VolumeNode{
			ID:    id,
			Label: fmt.Sprintf("Módulo inferior %d", i+1),
			Constraints: domain.VolumeConstraints{
				Width:  &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: ratios[i]},
				Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
			},
			Children: []domain.VolumeNode{},
			Features: []domain.Feature{
				{
					ID:     id + "-shelves",
					Type:   "shelf_set",
					Params: json.RawMessage(`{"count":2,"spacing":"equal"}`),
				},
			},
			Fronts: []domain.Front{
				{
					ID:     id + "-door",
					Type:   "door",
					Params: json.RawMessage(`{"hinge":"left","materialId":"melamine-white"}`),
				},
			},
		}
	}

	tvFronts := []domain.Front{}
	if tvFront.ID != "" {
		tvFronts = []domain.Front{tvFront}
	}

	return domain.FurnitureDefinition{
		ID:          fmt.Sprintf("ai-entertainment-%d", time.Now().Unix()),
		Name:        name,
		Description: description,
		Root: domain.VolumeNode{
			ID:    "root",
			Label: "Centro de entretenimiento",
			Constraints: domain.VolumeConstraints{
				Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFixed, Value: depth},
			},
			Split: &domain.VolumeSplit{
				Axis:   domain.SplitAxisY,
				Ratios: []float64{0.42, 0.58},
			},
			Children: []domain.VolumeNode{
				{
					ID:    "tv-bay",
					Label: "Nicho TV",
					Constraints: domain.VolumeConstraints{
						Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
						Height: &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: 0.42},
						Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
					},
					Children: []domain.VolumeNode{},
					Features: []domain.Feature{
						{
							ID:     "tv-space",
							Type:   "appliance_space",
							Params: json.RawMessage(`{"appliance":"tv","maxWidth":1800,"maxHeight":900}`),
						},
					},
					Fronts: tvFronts,
				},
				{
					ID:    "storage-base",
					Label: "Almacenaje inferior",
					Constraints: domain.VolumeConstraints{
						Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
						Height: &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: 0.58},
						Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
					},
					Split: &domain.VolumeSplit{
						Axis:   domain.SplitAxisX,
						Ratios: ratios,
					},
					Children: storageChildren,
					Features: []domain.Feature{},
					Fronts:   []domain.Front{},
				},
			},
			Features: []domain.Feature{},
			Fronts:   []domain.Front{},
			Adaptation: &domain.AdaptationRules{
				FollowFloor:   true,
				FollowCeiling: false,
			},
			Manufacturing: defaultManufacturingHints(),
		},
	}
}
