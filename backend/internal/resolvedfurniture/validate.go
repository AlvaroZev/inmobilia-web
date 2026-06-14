package resolvedfurniture

import (
	"fmt"

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

func ValidateResolvedFurniture(furniture domain.ResolvedFurniture) ValidationResult {
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	if furniture.ID == "" {
		result.add("id", "id is required")
	}
	if furniture.Name == "" {
		result.add("name", "name is required")
	}

	ids := map[string]string{}
	validateResolvedVolume(furniture.Root, "root", &result, ids, nil)

	return result
}

func validateResolvedVolume(
	volume domain.ResolvedVolume,
	prefix string,
	result *ValidationResult,
	ids map[string]string,
	parent *domain.ResolvedVolume,
) {
	if volume.ID == "" {
		result.add(prefix+".id", "id is required")
	} else if existing, ok := ids[volume.ID]; ok {
		result.add(prefix+".id", fmt.Sprintf(`duplicate id "%s" (also used in %s)`, volume.ID, existing))
	} else {
		ids[volume.ID] = prefix
	}

	validateDimensions(volume, prefix, result)

	if parent != nil && !VolumeContains(*parent, volume) {
		result.add(prefix, "volume exceeds parent bounds")
	}

	for i, child := range volume.Children {
		childPrefix := fmt.Sprintf("%s.children[%d]", prefix, i)
		validateResolvedVolume(child, childPrefix, result, ids, &volume)
	}

	for i := 0; i < len(volume.Children); i++ {
		for j := i + 1; j < len(volume.Children); j++ {
			if VolumeOverlap(volume.Children[i], volume.Children[j]) {
				result.add(
					fmt.Sprintf("%s.children[%d]", prefix, i),
					fmt.Sprintf(`overlaps with sibling "%s"`, volume.Children[j].ID),
				)
			}
		}
	}

	for i, feature := range volume.Features {
		field := fmt.Sprintf("%s.features[%d]", prefix, i)
		validateFeature(feature, field, result, ids)
		if !elementContainedInVolume(featureBounds(feature), volume) {
			result.add(field, "feature exceeds volume bounds")
		}
	}

	for i, front := range volume.Fronts {
		field := fmt.Sprintf("%s.fronts[%d]", prefix, i)
		validateFront(front, field, result, ids)
		if !elementContainedInVolume(frontBounds(front), volume) {
			result.add(field, "front exceeds volume bounds")
		}
	}
}

func validateDimensions(volume domain.ResolvedVolume, prefix string, result *ValidationResult) {
	if volume.Width <= geom.Epsilon {
		result.add(prefix+".width", "width must be greater than zero")
	}
	if volume.Height <= geom.Epsilon {
		result.add(prefix+".height", "height must be greater than zero")
	}
	if volume.Depth <= geom.Epsilon {
		result.add(prefix+".depth", "depth must be greater than zero")
	}
	if volume.Width < 0 || volume.Height < 0 || volume.Depth < 0 {
		result.add(prefix, "dimensions must be non-negative")
	}
}

func validateFeature(feature domain.ResolvedFeature, prefix string, result *ValidationResult, ids map[string]string) {
	if feature.ID == "" {
		result.add(prefix+".id", "id is required")
	} else if existing, ok := ids[feature.ID]; ok {
		result.add(prefix+".id", fmt.Sprintf(`duplicate id "%s" (also used in %s)`, feature.ID, existing))
	} else {
		ids[feature.ID] = prefix
	}
	if feature.Type == "" {
		result.add(prefix+".type", "type is required")
	}
	validateElementDimensions(feature.Width, feature.Height, feature.Depth, prefix, result)
}

func validateFront(front domain.ResolvedFront, prefix string, result *ValidationResult, ids map[string]string) {
	if front.ID == "" {
		result.add(prefix+".id", "id is required")
	} else if existing, ok := ids[front.ID]; ok {
		result.add(prefix+".id", fmt.Sprintf(`duplicate id "%s" (also used in %s)`, front.ID, existing))
	} else {
		ids[front.ID] = prefix
	}
	if front.Type == "" {
		result.add(prefix+".type", "type is required")
	}
	validateElementDimensions(front.Width, front.Height, front.Depth, prefix, result)
}

func validateElementDimensions(width, height, depth float64, prefix string, result *ValidationResult) {
	if width <= geom.Epsilon || height <= geom.Epsilon || depth <= geom.Epsilon {
		result.add(prefix, "dimensions must be greater than zero")
	}
}

func featureBounds(feature domain.ResolvedFeature) domain.BoundingBox3D {
	return domain.BoundingBox3D{
		Min: domain.Point3D{X: feature.X, Y: feature.Y, Z: feature.Z},
		Max: domain.Point3D{
			X: feature.X + feature.Width,
			Y: feature.Y + feature.Height,
			Z: feature.Z + feature.Depth,
		},
	}
}

func frontBounds(front domain.ResolvedFront) domain.BoundingBox3D {
	return domain.BoundingBox3D{
		Min: domain.Point3D{X: front.X, Y: front.Y, Z: front.Z},
		Max: domain.Point3D{
			X: front.X + front.Width,
			Y: front.Y + front.Height,
			Z: front.Z + front.Depth,
		},
	}
}

func elementContainedInVolume(bounds domain.BoundingBox3D, volume domain.ResolvedVolume) bool {
	vb := VolumeBounds(volume)
	bounds = geom.NormalizeBounds(bounds)

	return bounds.Min.X >= vb.Min.X-geom.Epsilon &&
		bounds.Min.Y >= vb.Min.Y-geom.Epsilon &&
		bounds.Min.Z >= vb.Min.Z-geom.Epsilon &&
		bounds.Max.X <= vb.Max.X+geom.Epsilon &&
		bounds.Max.Y <= vb.Max.Y+geom.Epsilon &&
		bounds.Max.Z <= vb.Max.Z+geom.Epsilon
}
