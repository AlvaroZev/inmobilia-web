package domain

import "encoding/json"

// ManufacturingModel (Layer 6) contains physical fabrication data.

type ManufacturingModel struct {
	FurnitureID string       `json:"furnitureId"`
	Parts       []Part       `json:"parts"`
	Hardware    []Hardware   `json:"hardware"`
	EdgeBanding []EdgeBanding `json:"edgeBanding"`
	Drilling    []Drilling   `json:"drilling"`
}

type Material struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Thickness float64 `json:"thickness"`
	Color     string  `json:"color,omitempty"`
}

type PartType string

const (
	PartLateral      PartType = "lateral"
	PartBase         PartType = "base"
	PartTop          PartType = "top"
	PartShelf        PartType = "shelf"
	PartDoor         PartType = "door"
	PartDivider      PartType = "divider"
	PartBack         PartType = "back"
	PartDrawerSide   PartType = "drawer_side"
	PartDrawerBottom PartType = "drawer_bottom"
	PartOther        PartType = "other"
)

type Part struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	VolumeID       string   `json:"volumeId"`
	Width          float64  `json:"width"`
	Height         float64  `json:"height"`
	Thickness      float64  `json:"thickness"`
	Material       Material `json:"material"`
	GrainDirection string   `json:"grainDirection,omitempty"`
}

type Hardware struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Quantity int             `json:"quantity"`
	PartIDs  []string        `json:"partIds,omitempty"`
	Params   json.RawMessage `json:"params,omitempty"`
}

type EdgeSide string

const (
	EdgeTop    EdgeSide = "top"
	EdgeBottom EdgeSide = "bottom"
	EdgeLeft   EdgeSide = "left"
	EdgeRight  EdgeSide = "right"
)

type EdgeBanding struct {
	PartID   string   `json:"partId"`
	Edge     EdgeSide `json:"edge"`
	Material string   `json:"material"`
	Length   float64  `json:"length"`
}

type Drilling struct {
	PartID    string  `json:"partId"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Diameter  float64 `json:"diameter"`
	Depth     float64 `json:"depth"`
	Type      string  `json:"type"`
}
