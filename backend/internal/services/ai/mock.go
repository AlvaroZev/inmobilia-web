package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

type MockParser struct{}

func NewMockParser() *MockParser {
	return &MockParser{}
}

func (p *MockParser) ParseFurniture(_ context.Context, description string, name string) (domain.FurnitureDefinition, error) {
	if strings.TrimSpace(description) == "" {
		return domain.FurnitureDefinition{}, ErrEmptyDescription
	}

	lower := strings.ToLower(description)
	switch {
	case isCloset(lower):
		return p.closetDefinition(description, name), nil
	case isEntertainmentCenter(lower):
		return p.entertainmentCenterDefinition(description, name), nil
	case isDesk(lower):
		return p.deskDefinition(description, name), nil
	default:
		return p.simpleCabinet(description, name), nil
	}
}

func (p *MockParser) closetDefinition(description, name string) domain.FurnitureDefinition {
	if name == "" {
		name = "Closet empotrado"
	}

	lower := strings.ToLower(description)
	bodyCount := parseBodyCount(description)
	ratios := equalRatios(bodyCount)
	withDrawers := strings.Contains(lower, "cajon")

	children := make([]domain.VolumeNode, bodyCount)
	for i := range bodyCount {
		id := fmt.Sprintf("body-%d", i+1)
		label := fmt.Sprintf("Cuerpo %d", i+1)
		withRod := i == 0 && !(withDrawers && i == bodyCount-1)
		withDrawersOnBody := withDrawers && i == bodyCount-1
		children[i] = p.closetBody(id, label, ratios[i], withRod, withDrawersOnBody)
	}

	return domain.FurnitureDefinition{
		ID:          fmt.Sprintf("ai-closet-%d", time.Now().Unix()),
		Name:        name,
		Description: description,
		Root: domain.VolumeNode{
			ID:    "root",
			Label: "Cuerpo principal",
			Constraints: domain.VolumeConstraints{
				Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFixed, Value: 600},
			},
			Split: &domain.VolumeSplit{
				Axis:   domain.SplitAxisX,
				Ratios: ratios,
			},
			Children: children,
			Features: []domain.Feature{},
			Fronts:   []domain.Front{},
			Adaptation: &domain.AdaptationRules{
				FollowFloor:        true,
				FollowCeiling:      true,
				CompensateSkirting: true,
			},
			Manufacturing: defaultManufacturingHints(),
		},
	}
}

func (p *MockParser) closetBody(id, label string, ratio float64, withRod, withDrawers bool) domain.VolumeNode {
	node := domain.VolumeNode{
		ID:    id,
		Label: label,
		Constraints: domain.VolumeConstraints{
			Width:  &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: ratio},
			Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
			Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
		},
		Children: []domain.VolumeNode{},
		Features: []domain.Feature{
			{
				ID:     id + "-shelves",
				Type:   "shelf_set",
				Params: json.RawMessage(`{"count":4,"spacing":"equal"}`),
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

	if withRod {
		node.Features = append(node.Features, domain.Feature{
			ID:     id + "-rod",
			Type:   "hanger_rod",
			Params: json.RawMessage(`{"heightFromTop":1800}`),
		})
	}

	if withDrawers {
		node.Split = &domain.VolumeSplit{
			Axis:   domain.SplitAxisY,
			Ratios: []float64{0.7, 0.3},
		}
		node.Children = []domain.VolumeNode{
			{
				ID:    id + "-upper",
				Label: "Repisas superiores",
				Constraints: domain.VolumeConstraints{
					Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
					Height: &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: 0.7},
					Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				},
				Children: []domain.VolumeNode{},
				Features: []domain.Feature{
					{
						ID:     id + "-shelves-upper",
						Type:   "shelf_set",
						Params: json.RawMessage(`{"count":3,"spacing":"equal"}`),
					},
				},
				Fronts: []domain.Front{
					{
						ID:     id + "-door-upper",
						Type:   "door",
						Params: json.RawMessage(`{"hinge":"right","materialId":"melamine-white"}`),
					},
				},
			},
			{
				ID:    id + "-drawers",
				Label: "Cajonera",
				Constraints: domain.VolumeConstraints{
					Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
					Height: &domain.DimensionConstraint{Mode: domain.DimensionRatio, Value: 0.3},
					Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				},
				Children: []domain.VolumeNode{},
				Features: []domain.Feature{
					{
						ID:     id + "-drawers-stack",
						Type:   "drawer_stack",
						Params: json.RawMessage(`{"count":3,"runner":"soft-close"}`),
					},
				},
				Fronts: []domain.Front{
					{
						ID:     id + "-drawer-fronts",
						Type:   "drawer_front",
						Params: json.RawMessage(`{"materialId":"melamine-white"}`),
					},
				},
			},
		}
		node.Features = []domain.Feature{}
		node.Fronts = []domain.Front{}
	}

	return node
}

func (p *MockParser) simpleCabinet(description, name string) domain.FurnitureDefinition {
	if name == "" {
		name = "Mueble a medida"
	}

	depth := 450.0
	if strings.Contains(strings.ToLower(description), "600") {
		depth = 600
	}

	return domain.FurnitureDefinition{
		ID:          fmt.Sprintf("ai-cabinet-%d", time.Now().Unix()),
		Name:        name,
		Description: description,
		Root: domain.VolumeNode{
			ID:    "root",
			Label: "Cuerpo",
			Constraints: domain.VolumeConstraints{
				Width:  &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Height: &domain.DimensionConstraint{Mode: domain.DimensionFill},
				Depth:  &domain.DimensionConstraint{Mode: domain.DimensionFixed, Value: depth},
			},
			Children: []domain.VolumeNode{},
			Features: []domain.Feature{
				{
					ID:     "root-shelves",
					Type:   "shelf_set",
					Params: json.RawMessage(`{"count":3,"spacing":"equal"}`),
				},
			},
			Fronts: []domain.Front{},
			Adaptation: &domain.AdaptationRules{
				FollowFloor:   true,
				FollowCeiling: true,
			},
			Manufacturing: &domain.ManufacturingHints{
				MaterialID: "melamine-white-18",
				BackPanel:  true,
			},
		},
	}
}
