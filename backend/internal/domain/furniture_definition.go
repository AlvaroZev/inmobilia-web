package domain

import "encoding/json"

// FurnitureDefinition (Layer 2) is design intent only — no final geometry or parts.

type FurnitureDefinition struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Root        VolumeNode      `json:"root"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

type VolumeNode struct {
	ID             string              `json:"id"`
	Label          string              `json:"label,omitempty"`
	Constraints    VolumeConstraints   `json:"constraints"`
	Split          *VolumeSplit        `json:"split,omitempty"`
	Children       []VolumeNode        `json:"children"`
	Features       []Feature           `json:"features"`
	Fronts         []Front             `json:"fronts"`
	Adaptation     *AdaptationRules    `json:"adaptation,omitempty"`
	Manufacturing  *ManufacturingHints `json:"manufacturing,omitempty"`
}

type VolumeConstraints struct {
	Width  *DimensionConstraint `json:"width,omitempty"`
	Height *DimensionConstraint `json:"height,omitempty"`
	Depth  *DimensionConstraint `json:"depth,omitempty"`
}

type DimensionMode string

const (
	DimensionFixed DimensionMode = "fixed"
	DimensionRatio DimensionMode = "ratio"
	DimensionFill  DimensionMode = "fill"
	DimensionMin   DimensionMode = "min"
	DimensionMax   DimensionMode = "max"
)

type DimensionConstraint struct {
	Mode  DimensionMode `json:"mode"`
	Value float64       `json:"value,omitempty"`
}

type SplitAxis string

const (
	SplitAxisX SplitAxis = "x"
	SplitAxisY SplitAxis = "y"
	SplitAxisZ SplitAxis = "z"
)

type VolumeSplit struct {
	Axis   SplitAxis `json:"axis"`
	Ratios []float64 `json:"ratios,omitempty"`
	Fixed  []float64 `json:"fixed,omitempty"`
}

// Feature is extensible — type is a string, not a closed enum.
type Feature struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

type Front struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

type AdaptationRules struct {
	FollowFloor         bool            `json:"followFloor,omitempty"`
	FollowCeiling       bool            `json:"followCeiling,omitempty"`
	FollowWall          bool            `json:"followWall,omitempty"`
	CompensateSkirting  bool            `json:"compensateSkirting,omitempty"`
	Params              json.RawMessage `json:"params,omitempty"`
}

type ManufacturingHints struct {
	MaterialID   string          `json:"materialId,omitempty"`
	EdgeBanding  string          `json:"edgeBanding,omitempty"`
	BackPanel    bool            `json:"backPanel,omitempty"`
	Params       json.RawMessage `json:"params,omitempty"`
}
