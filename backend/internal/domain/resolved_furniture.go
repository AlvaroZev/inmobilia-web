package domain

import "encoding/json"

// ResolvedFurniture (Layer 5) contains real computed dimensions.

type ResolvedFurniture struct {
	ID   string         `json:"id"`
	Name string         `json:"name"`
	Root ResolvedVolume `json:"root"`
}

type ResolvedVolume struct {
	ID         string           `json:"id"`
	Label      string           `json:"label,omitempty"`
	X          float64          `json:"x"`
	Y          float64          `json:"y"`
	Z          float64          `json:"z"`
	Width      float64          `json:"width"`
	Height     float64          `json:"height"`
	Depth      float64          `json:"depth"`
	Children   []ResolvedVolume `json:"children"`
	Features   []ResolvedFeature `json:"features"`
	Fronts     []ResolvedFront  `json:"fronts"`
	MaterialID string           `json:"materialId,omitempty"`
}

type ResolvedFeature struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	X      float64         `json:"x"`
	Y      float64         `json:"y"`
	Z      float64         `json:"z"`
	Width  float64         `json:"width"`
	Height float64         `json:"height"`
	Depth  float64         `json:"depth"`
	Params json.RawMessage `json:"params"`
}

type ResolvedFront struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	X      float64         `json:"x"`
	Y      float64         `json:"y"`
	Z      float64         `json:"z"`
	Width  float64         `json:"width"`
	Height float64         `json:"height"`
	Depth  float64         `json:"depth"`
	Params json.RawMessage `json:"params"`
}
