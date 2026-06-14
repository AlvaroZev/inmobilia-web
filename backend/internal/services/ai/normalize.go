package ai

import (
	"encoding/json"
	"fmt"
)

func normalizeAIJSON(raw []byte) ([]byte, error) {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	root, ok := data["root"].(map[string]any)
	if !ok {
		repairFurnitureDocument(data)
		return json.Marshal(data)
	}

	normalizeVolumeNode(root, "root")
	repairFurnitureDocument(data)

	return json.Marshal(data)
}

func normalizeVolumeNode(node map[string]any, prefix string) {
	switch children := node["children"].(type) {
	case []any:
		for i, child := range children {
			childMap, ok := child.(map[string]any)
			if !ok {
				continue
			}
			childPrefix := fmt.Sprintf("%s-c%d", prefix, i)
			normalizeVolumeNode(childMap, childPrefix)
			children[i] = childMap
		}
		node["children"] = children
	default:
		node["children"] = []any{}
	}

	node["features"] = normalizeFeatures(node["features"], prefix)
	node["fronts"] = normalizeFronts(node["fronts"], prefix)
}

func normalizeFeatures(raw any, prefix string) []any {
	items, ok := raw.([]any)
	if !ok || len(items) == 0 {
		return []any{}
	}

	out := make([]any, 0, len(items))
	for i, item := range items {
		id := fmt.Sprintf("%s-f%d", prefix, i)
		switch value := item.(type) {
		case string:
			out = append(out, featureFromString(value, id))
		case map[string]any:
			out = append(out, normalizeFeatureObject(value, id))
		}
	}
	return out
}

func normalizeFronts(raw any, prefix string) []any {
	items, ok := raw.([]any)
	if !ok || len(items) == 0 {
		return []any{}
	}

	out := make([]any, 0, len(items))
	for i, item := range items {
		id := fmt.Sprintf("%s-fr%d", prefix, i)
		switch value := item.(type) {
		case string:
			out = append(out, frontFromString(value, id))
		case map[string]any:
			out = append(out, normalizeFrontObject(value, id))
		}
	}
	return out
}

func normalizeFeatureObject(raw map[string]any, fallbackID string) map[string]any {
	feature := map[string]any{
		"id":     stringField(raw, "id", fallbackID),
		"type":   stringField(raw, "type", "custom"),
		"params": objectField(raw, "params"),
	}
	if feature["type"] == "custom" {
		if typeName, ok := raw["type"].(string); ok && typeName != "" {
			feature["type"] = typeName
		}
	}
	if params, ok := feature["params"].(map[string]any); !ok || len(params) == 0 {
		feature["params"] = defaultFeatureParams(feature["type"].(string))
	}
	return feature
}

func normalizeFrontObject(raw map[string]any, fallbackID string) map[string]any {
	front := map[string]any{
		"id":     stringField(raw, "id", fallbackID),
		"type":   stringField(raw, "type", "door"),
		"params": objectField(raw, "params"),
	}
	if params, ok := front["params"].(map[string]any); !ok || len(params) == 0 {
		front["params"] = defaultFrontParams(front["type"].(string))
	}
	return front
}

func featureFromString(typeName, id string) map[string]any {
	return map[string]any{
		"id":     id,
		"type":   typeName,
		"params": defaultFeatureParams(typeName),
	}
}

func frontFromString(typeName, id string) map[string]any {
	return map[string]any{
		"id":     id,
		"type":   typeName,
		"params": defaultFrontParams(typeName),
	}
}

func defaultFeatureParams(typeName string) map[string]any {
	switch typeName {
	case "shelf_set":
		return map[string]any{"count": 4, "spacing": "equal"}
	case "drawer_stack":
		return map[string]any{"count": 3, "runner": "soft-close"}
	case "hanger_rod":
		return map[string]any{"heightFromTop": 1800}
	case "appliance_space":
		return map[string]any{"appliance": "tv"}
	default:
		return map[string]any{}
	}
}

func defaultFrontParams(typeName string) map[string]any {
	switch typeName {
	case "door":
		return map[string]any{"hinge": "left", "materialId": "melamine-white"}
	case "drawer_front":
		return map[string]any{"materialId": "melamine-white"}
	default:
		return map[string]any{}
	}
}

func stringField(raw map[string]any, key, fallback string) string {
	if value, ok := raw[key].(string); ok && value != "" {
		return value
	}
	return fallback
}

func objectField(raw map[string]any, key string) map[string]any {
	switch value := raw[key].(type) {
	case map[string]any:
		return value
	case string:
		var parsed map[string]any
		if err := json.Unmarshal([]byte(value), &parsed); err == nil {
			return parsed
		}
	}
	return map[string]any{}
}
