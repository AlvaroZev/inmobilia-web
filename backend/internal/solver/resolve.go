package solver

import (
	"encoding/json"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
)

func resolveVolumeTree(node domain.VolumeNode, slot box, splitAxis *domain.SplitAxis) (domain.ResolvedVolume, error) {
	resolvedBox, err := resolveNodeBox(node, slot, splitAxis)
	if err != nil {
		return domain.ResolvedVolume{}, err
	}

	resolved := domain.ResolvedVolume{
		ID:         node.ID,
		Label:      node.Label,
		X:          resolvedBox.X,
		Y:          resolvedBox.Y,
		Z:          resolvedBox.Z,
		Width:      resolvedBox.Width,
		Height:     resolvedBox.Height,
		Depth:      resolvedBox.Depth,
		Children:   []domain.ResolvedVolume{},
		Features:   resolveFeatures(node.Features, resolvedBox),
		Fronts:     resolveFronts(node.Fronts, resolvedBox),
		MaterialID: materialID(node),
	}

	if len(node.Children) == 0 {
		return resolved, nil
	}

	if node.Split == nil {
		return domain.ResolvedVolume{}, ErrChildrenWithoutSplit
	}

	childSlots, err := allocateChildBoxes(node, resolvedBox)
	if err != nil {
		return domain.ResolvedVolume{}, err
	}

	axis := node.Split.Axis
	for i, child := range node.Children {
		childResolved, err := resolveVolumeTree(child, childSlots[i], &axis)
		if err != nil {
			return domain.ResolvedVolume{}, err
		}
		resolved.Children = append(resolved.Children, childResolved)
	}

	return resolved, nil
}

func resolveNodeBox(node domain.VolumeNode, slot box, lockedAxis *domain.SplitAxis) (box, error) {
	result := slot

	width, err := resolveAxisDimension(slot.Width, node.Constraints.Width, lockedAxis, domain.SplitAxisX)
	if err != nil {
		return box{}, err
	}
	height, err := resolveAxisDimension(slot.Height, node.Constraints.Height, lockedAxis, domain.SplitAxisY)
	if err != nil {
		return box{}, err
	}
	depth, err := resolveAxisDimension(slot.Depth, node.Constraints.Depth, lockedAxis, domain.SplitAxisZ)
	if err != nil {
		return box{}, err
	}

	result.Width = width
	result.Height = height
	result.Depth = depth

	return result, nil
}

func resolveAxisDimension(
	available float64,
	constraint *domain.DimensionConstraint,
	lockedAxis *domain.SplitAxis,
	axis domain.SplitAxis,
) (float64, error) {
	if lockedAxis != nil && *lockedAxis == axis {
		return available, nil
	}
	return resolveDimension(available, constraint)
}

func resolveDimension(available float64, constraint *domain.DimensionConstraint) (float64, error) {
	if constraint == nil {
		return available, nil
	}

	switch constraint.Mode {
	case domain.DimensionFill:
		return available, nil
	case domain.DimensionFixed:
		if constraint.Value > available+geom.Epsilon {
			return 0, ErrDimensionExceedsSpace
		}
		return constraint.Value, nil
	case domain.DimensionRatio:
		return available * constraint.Value, nil
	case domain.DimensionMin:
		if constraint.Value > available+geom.Epsilon {
			return 0, ErrDimensionExceedsSpace
		}
		return available, nil
	case domain.DimensionMax:
		value := constraint.Value
		if value > available {
			value = available
		}
		return value, nil
	default:
		return 0, ErrUnknownConstraintMode
	}
}

func allocateChildBoxes(parent domain.VolumeNode, parentBox box) ([]box, error) {
	if parent.Split == nil {
		return nil, ErrChildrenWithoutSplit
	}

	n := len(parent.Children)
	slots := make([]box, n)
	cursor := parentBox

	for i := 0; i < n; i++ {
		seg, err := splitSegmentSize(parent.Split, parentBox, i)
		if err != nil {
			return nil, err
		}

		child := parentBox
		switch parent.Split.Axis {
		case domain.SplitAxisX:
			child.X = cursor.X
			child.Width = seg
			cursor.X += seg
		case domain.SplitAxisY:
			child.Y = cursor.Y
			child.Height = seg
			cursor.Y += seg
		case domain.SplitAxisZ:
			child.Z = cursor.Z
			child.Depth = seg
			cursor.Z += seg
		default:
			return nil, ErrInvalidSplitAxis
		}

		slots[i] = child
	}

	return slots, nil
}

func splitSegmentSize(split *domain.VolumeSplit, parentBox box, index int) (float64, error) {
	total := parentBox.dimension(split.Axis)

	if len(split.Ratios) > 0 {
		if index >= len(split.Ratios) {
			return 0, ErrSplitChildMismatch
		}
		return total * split.Ratios[index], nil
	}

	if len(split.Fixed) > 0 {
		if index >= len(split.Fixed) {
			return 0, ErrSplitChildMismatch
		}
		if split.Fixed[index] > total+geom.Epsilon {
			return 0, ErrDimensionExceedsSpace
		}
		return split.Fixed[index], nil
	}

	return 0, ErrSplitUndefined
}

func resolveFeatures(features []domain.Feature, b box) []domain.ResolvedFeature {
	resolved := make([]domain.ResolvedFeature, len(features))
	for i, feature := range features {
		resolved[i] = domain.ResolvedFeature{
			ID:     feature.ID,
			Type:   feature.Type,
			X:      b.X,
			Y:      b.Y,
			Z:      b.Z,
			Width:  b.Width,
			Height: b.Height,
			Depth:  b.Depth,
			Params: normalizeParams(feature.Params),
		}
	}
	return resolved
}

func resolveFronts(fronts []domain.Front, b box) []domain.ResolvedFront {
	resolved := make([]domain.ResolvedFront, len(fronts))
	for i, front := range fronts {
		resolved[i] = domain.ResolvedFront{
			ID:     front.ID,
			Type:   front.Type,
			X:      b.X,
			Y:      b.Y,
			Z:      b.Z,
			Width:  b.Width,
			Height: b.Height,
			Depth:  b.Depth,
			Params: normalizeParams(front.Params),
		}
	}
	return resolved
}

func normalizeParams(params json.RawMessage) json.RawMessage {
	if len(params) == 0 {
		return json.RawMessage("{}")
	}
	return params
}

func materialID(node domain.VolumeNode) string {
	if node.Manufacturing == nil {
		return ""
	}
	return node.Manufacturing.MaterialID
}
