package volumetree

import (
	"fmt"
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/pkg/geom"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

func (r *ValidationResult) add(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
	r.Valid = false
}

func ValidateFurnitureDefinition(furniture domain.FurnitureDefinition) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	if furniture.ID == "" {
		result.add("id", "id is required")
	}
	if furniture.Name == "" {
		result.add("name", "name is required")
	}

	ids := map[string]string{}
	validateVolumeNode(furniture.Root, "root", &result, ids)

	return result
}

func validateVolumeNode(node domain.VolumeNode, prefix string, result *ValidationResult, ids map[string]string) {
	if node.ID == "" {
		result.add(prefix+".id", "id is required")
	} else if existing, ok := ids[node.ID]; ok {
		result.add(prefix+".id", fmt.Sprintf(`duplicate id "%s" (also used in %s)`, node.ID, existing))
	} else {
		ids[node.ID] = prefix
	}

	validateConstraints(node.Constraints, prefix, result)

	hasChildren := len(node.Children) > 0
	hasSplit := node.Split != nil

	if hasChildren && !hasSplit {
		result.add(prefix+".children", "nodes with children must define a split")
	}
	if hasSplit && !hasChildren {
		result.add(prefix+".split", "split requires at least one child")
	}
	if hasSplit && hasChildren {
		validateSplit(*node.Split, len(node.Children), prefix, result)
	}

	for i, feature := range node.Features {
		field := fmt.Sprintf("%s.features[%d].id", prefix, i)
		validateFeature(feature, fmt.Sprintf("%s.features[%d]", prefix, i), result)
		if feature.ID != "" {
			if existing, ok := ids[feature.ID]; ok {
				result.add(field, fmt.Sprintf(`duplicate id "%s" (also used in %s)`, feature.ID, existing))
			} else {
				ids[feature.ID] = field
			}
		}
	}

	for i, front := range node.Fronts {
		field := fmt.Sprintf("%s.fronts[%d].id", prefix, i)
		validateFront(front, fmt.Sprintf("%s.fronts[%d]", prefix, i), result)
		if front.ID != "" {
			if existing, ok := ids[front.ID]; ok {
				result.add(field, fmt.Sprintf(`duplicate id "%s" (also used in %s)`, front.ID, existing))
			} else {
				ids[front.ID] = field
			}
		}
	}

	for i, child := range node.Children {
		childPrefix := fmt.Sprintf("%s.children[%d]", prefix, i)
		validateVolumeNode(child, childPrefix, result, ids)
		validateChildSplitAlignment(node, child, i, prefix, result)
	}
}

func validateConstraints(constraints domain.VolumeConstraints, prefix string, result *ValidationResult) {
	if constraints.Width != nil {
		validateDimensionConstraint(*constraints.Width, prefix+".constraints.width", result)
	}
	if constraints.Height != nil {
		validateDimensionConstraint(*constraints.Height, prefix+".constraints.height", result)
	}
	if constraints.Depth != nil {
		validateDimensionConstraint(*constraints.Depth, prefix+".constraints.depth", result)
	}
}

func validateDimensionConstraint(constraint domain.DimensionConstraint, field string, result *ValidationResult) {
	switch constraint.Mode {
	case domain.DimensionFixed, domain.DimensionMin, domain.DimensionMax:
		if constraint.Value <= geom.Epsilon {
			result.add(field, fmt.Sprintf(`mode "%s" requires a positive value`, constraint.Mode))
		}
	case domain.DimensionRatio:
		if constraint.Value <= geom.Epsilon || constraint.Value > 1+geom.Epsilon {
			result.add(field, `mode "ratio" requires a value between 0 and 1`)
		}
	case domain.DimensionFill:
	default:
		result.add(field, fmt.Sprintf(`unknown constraint mode "%s"`, constraint.Mode))
	}
}

func validateSplit(split domain.VolumeSplit, childCount int, prefix string, result *ValidationResult) {
	switch split.Axis {
	case domain.SplitAxisX, domain.SplitAxisY, domain.SplitAxisZ:
	default:
		result.add(prefix+".split.axis", "axis must be one of: x, y, z")
	}

	hasRatios := len(split.Ratios) > 0
	hasFixed := len(split.Fixed) > 0

	if !hasRatios && !hasFixed {
		result.add(prefix+".split", "split requires ratios or fixed sizes")
		return
	}
	if hasRatios && hasFixed {
		result.add(prefix+".split", "split cannot define both ratios and fixed sizes")
	}

	if hasRatios {
		if len(split.Ratios) != childCount {
			result.add(prefix+".split.ratios",
				fmt.Sprintf("ratio count (%d) must match children count (%d)", len(split.Ratios), childCount))
		}
		for i, ratio := range split.Ratios {
			if ratio <= geom.Epsilon {
				result.add(fmt.Sprintf("%s.split.ratios[%d]", prefix, i), "ratio must be greater than zero")
			}
		}
		total := SumSplitRatios(split)
		if math.Abs(total-1) > 0.01 {
			result.add(prefix+".split.ratios", fmt.Sprintf("ratios must sum to 1 (got %v)", total))
		}
	}

	if hasFixed {
		if len(split.Fixed) != childCount {
			result.add(prefix+".split.fixed",
				fmt.Sprintf("fixed count (%d) must match children count (%d)", len(split.Fixed), childCount))
		}
		for i, fixed := range split.Fixed {
			if fixed <= geom.Epsilon {
				result.add(fmt.Sprintf("%s.split.fixed[%d]", prefix, i), "fixed size must be greater than zero")
			}
		}
	}
}

func validateChildSplitAlignment(
	parent domain.VolumeNode,
	child domain.VolumeNode,
	childIndex int,
	prefix string,
	result *ValidationResult,
) {
	if parent.Split == nil {
		return
	}

	dim := AxisToDimension(parent.Split.Axis)
	var childConstraint *domain.DimensionConstraint

	switch dim {
	case "width":
		childConstraint = child.Constraints.Width
	case "height":
		childConstraint = child.Constraints.Height
	case "depth":
		childConstraint = child.Constraints.Depth
	}

	childField := fmt.Sprintf("%s.children[%d].constraints.%s", prefix, childIndex, dim)

	if childConstraint == nil {
		result.add(childField,
			fmt.Sprintf(`child must define a %s constraint matching parent split on axis "%s"`, dim, parent.Split.Axis))
		return
	}

	if len(parent.Split.Ratios) > 0 {
		expected := parent.Split.Ratios[childIndex]
		if childConstraint.Mode == domain.DimensionRatio && math.Abs(childConstraint.Value-expected) > 0.01 {
			result.add(childField,
				fmt.Sprintf("ratio value %v does not match parent split ratio %v", childConstraint.Value, expected))
		}
		if childConstraint.Mode == domain.DimensionFixed {
			result.add(childField, "child cannot use fixed constraint when parent splits by ratios")
		}
	}

	if len(parent.Split.Fixed) > 0 {
		expected := parent.Split.Fixed[childIndex]
		if childConstraint.Mode == domain.DimensionFixed && math.Abs(childConstraint.Value-expected) > geom.Epsilon {
			result.add(childField,
				fmt.Sprintf("fixed value %v does not match parent split fixed %v", childConstraint.Value, expected))
		}
		if childConstraint.Mode == domain.DimensionRatio {
			result.add(childField, "child cannot use ratio constraint when parent splits by fixed sizes")
		}
	}
}

func validateFeature(feature domain.Feature, prefix string, result *ValidationResult) {
	if feature.ID == "" {
		result.add(prefix+".id", "id is required")
	}
	if feature.Type == "" {
		result.add(prefix+".type", "type is required")
	}
}

func validateFront(front domain.Front, prefix string, result *ValidationResult) {
	if front.ID == "" {
		result.add(prefix+".id", "id is required")
	}
	if front.Type == "" {
		result.add(prefix+".type", "type is required")
	}
}
