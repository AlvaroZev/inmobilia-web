package ai

import (
	"fmt"
	"math"
	"time"
)

func repairFurnitureDocument(data map[string]any) {
	ensureFurnitureMeta(data)

	root, ok := data["root"].(map[string]any)
	if !ok {
		return
	}

	usedIDs := map[string]int{}
	repairVolumeTree(root, "root", usedIDs)
	data["root"] = root
}

func ensureFurnitureMeta(data map[string]any) {
	if stringField(data, "id", "") == "" {
		data["id"] = fmt.Sprintf("ai-%d", time.Now().Unix())
	}
	if stringField(data, "name", "") == "" {
		data["name"] = "Mueble a medida"
	}
	if _, ok := data["description"]; !ok {
		data["description"] = ""
	}
}

func repairVolumeTree(node map[string]any, fallbackID string, usedIDs map[string]int) {
	node["id"] = uniqueID(stringField(node, "id", fallbackID), usedIDs)
	ensureNodeConstraints(node)

	children, _ := node["children"].([]any)
	split, hasSplit := node["split"].(map[string]any)

	if len(children) > 0 {
		if !hasSplit {
			node["split"] = map[string]any{
				"axis":   "x",
				"ratios": equalRatioValues(len(children)),
			}
			split, _ = node["split"].(map[string]any)
			hasSplit = true
		}

		if hasSplit && split != nil {
			repairSplit(split, len(children))
			axis := stringField(split, "axis", "x")
			ratios := toFloatSlice(split["ratios"])

			for i, child := range children {
				childMap, ok := child.(map[string]any)
				if !ok {
					continue
				}
				childPrefix := fmt.Sprintf("%s-c%d", fallbackID, i)
				repairChildConstraints(childMap, axis, ratios, i)
				repairVolumeTree(childMap, childPrefix, usedIDs)
				children[i] = childMap
			}
			node["children"] = children
		}
	} else if hasSplit {
		delete(node, "split")
	}

	repairFeatureFrontIDs(node, fallbackID, usedIDs)
}

func repairFeatureFrontIDs(node map[string]any, prefix string, usedIDs map[string]int) {
	if features, ok := node["features"].([]any); ok {
		for i, item := range features {
			feature, ok := item.(map[string]any)
			if !ok {
				continue
			}
			feature["id"] = uniqueID(stringField(feature, "id", fmt.Sprintf("%s-f%d", prefix, i)), usedIDs)
			if stringField(feature, "type", "") == "" {
				feature["type"] = "custom"
			}
			features[i] = feature
		}
		node["features"] = features
	}

	if fronts, ok := node["fronts"].([]any); ok {
		for i, item := range fronts {
			front, ok := item.(map[string]any)
			if !ok {
				continue
			}
			front["id"] = uniqueID(stringField(front, "id", fmt.Sprintf("%s-fr%d", prefix, i)), usedIDs)
			if stringField(front, "type", "") == "" {
				front["type"] = "door"
			}
			fronts[i] = front
		}
		node["fronts"] = fronts
	}
}

func ensureNodeConstraints(node map[string]any) {
	constraints, ok := node["constraints"].(map[string]any)
	if !ok {
		constraints = map[string]any{}
		node["constraints"] = constraints
	}

	for _, dim := range []string{"width", "height", "depth"} {
		if _, exists := constraints[dim]; !exists {
			constraints[dim] = map[string]any{"mode": "fill"}
		}
	}

	repairConstraintObjects(constraints)
}

func repairConstraintObjects(constraints map[string]any) {
	for _, dim := range []string{"width", "height", "depth"} {
		raw, ok := constraints[dim].(map[string]any)
		if !ok {
			constraints[dim] = map[string]any{"mode": "fill"}
			continue
		}

		mode := stringField(raw, "mode", "fill")
		switch mode {
		case "fixed", "ratio", "fill", "min", "max":
			raw["mode"] = mode
		default:
			raw["mode"] = "fill"
			mode = "fill"
		}

		if mode == "fixed" || mode == "ratio" || mode == "min" || mode == "max" {
			if _, hasValue := raw["value"]; !hasValue {
				if mode == "ratio" {
					raw["value"] = 0.5
				} else {
					raw["value"] = 600.0
				}
			}
		}

		constraints[dim] = raw
	}
}

func repairSplit(split map[string]any, childCount int) {
	axis := stringField(split, "axis", "x")
	if axis != "x" && axis != "y" && axis != "z" {
		axis = "x"
	}
	split["axis"] = axis

	ratios := toFloatSlice(split["ratios"])
	if len(ratios) != childCount || len(ratios) == 0 {
		ratios = equalRatioFloats(childCount)
	}
	split["ratios"] = normalizeRatioValues(ratios)
	delete(split, "fixed")
}

func repairChildConstraints(child map[string]any, axis string, ratios []float64, index int) {
	constraints, ok := child["constraints"].(map[string]any)
	if !ok {
		constraints = map[string]any{}
		child["constraints"] = constraints
	}

	splitDim := axisToDimension(axis)
	ratioValue := 1.0
	if index < len(ratios) {
		ratioValue = ratios[index]
	}
	constraints[splitDim] = map[string]any{"mode": "ratio", "value": ratioValue}

	for _, dim := range []string{"width", "height", "depth"} {
		if dim == splitDim {
			continue
		}
		if _, exists := constraints[dim]; !exists {
			constraints[dim] = map[string]any{"mode": "fill"}
		}
	}

	repairConstraintObjects(constraints)
}

func axisToDimension(axis string) string {
	switch axis {
	case "y":
		return "height"
	case "z":
		return "depth"
	default:
		return "width"
	}
}

func equalRatioValues(count int) []any {
	if count <= 0 {
		return []any{}
	}
	share := 1.0 / float64(count)
	values := make([]any, count)
	for i := range values {
		values[i] = share
	}
	return values
}

func equalRatioFloats(count int) []float64 {
	if count <= 0 {
		return nil
	}
	share := 1.0 / float64(count)
	values := make([]float64, count)
	for i := range values {
		values[i] = share
	}
	return values
}

func normalizeRatioValues(ratios []float64) []any {
	if len(ratios) == 0 {
		return []any{}
	}

	total := 0.0
	for _, ratio := range ratios {
		if ratio > 0 {
			total += ratio
		}
	}
	if total <= math.SmallestNonzeroFloat64 {
		return equalRatioValues(len(ratios))
	}

	out := make([]any, len(ratios))
	for i, ratio := range ratios {
		if ratio <= 0 {
			ratio = total / float64(len(ratios))
		}
		out[i] = ratio / total
	}
	return out
}

func toFloatSlice(raw any) []float64 {
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]float64, 0, len(items))
	for _, item := range items {
		switch value := item.(type) {
		case float64:
			out = append(out, value)
		case int:
			out = append(out, float64(value))
		}
	}
	return out
}

func uniqueID(id string, used map[string]int) string {
	if id == "" {
		id = "node"
	}
	count := used[id]
	if count == 0 {
		used[id] = 1
		return id
	}
	count++
	used[id] = count
	newID := fmt.Sprintf("%s-%d", id, count)
	used[newID] = 1
	return newID
}
