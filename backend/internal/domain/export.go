package domain

import "time"

type BillOfMaterials struct {
	FurnitureID   string             `json:"furnitureId"`
	FurnitureName string             `json:"furnitureName,omitempty"`
	GeneratedAt   time.Time          `json:"generatedAt"`
	Parts         []BOMPartLine      `json:"parts"`
	Hardware      []BOMHardwareLine  `json:"hardware"`
	EdgeBanding   []BOMEdgeLine      `json:"edgeBanding"`
	Cost          *CostResult        `json:"cost,omitempty"`
	Summary       BOMSummary         `json:"summary"`
}

type BOMPartLine struct {
	PartID         string  `json:"partId"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	VolumeID       string  `json:"volumeId"`
	Width          float64 `json:"width"`
	Height         float64 `json:"height"`
	Thickness      float64 `json:"thickness"`
	MaterialID     string  `json:"materialId"`
	MaterialName   string  `json:"materialName"`
	GrainDirection string  `json:"grainDirection,omitempty"`
	AreaM2         float64 `json:"areaM2"`
}

type BOMHardwareLine struct {
	HardwareType string  `json:"hardwareType"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	UnitCost     float64 `json:"unitCost,omitempty"`
	Total        float64 `json:"total,omitempty"`
}

type BOMEdgeLine struct {
	Material   string  `json:"material"`
	TotalLengthM float64 `json:"totalLengthM"`
}

type BOMSummary struct {
	PartCount      int     `json:"partCount"`
	HardwareCount  int     `json:"hardwareCount"`
	TotalBoardM2   float64 `json:"totalBoardM2"`
	TotalEdgeM     float64 `json:"totalEdgeM"`
}

type CutPlan struct {
	FurnitureID   string      `json:"furnitureId"`
	FurnitureName string      `json:"furnitureName,omitempty"`
	GeneratedAt   time.Time   `json:"generatedAt"`
	Sheets        []CutSheet  `json:"sheets"`
}

type CutSheet struct {
	MaterialID   string        `json:"materialId"`
	MaterialName string        `json:"materialName"`
	Thickness    float64       `json:"thickness"`
	Parts        []CutPartLine `json:"parts"`
	TotalAreaM2  float64       `json:"totalAreaM2"`
}

type CutPartLine struct {
	PartID    string  `json:"partId"`
	Name      string  `json:"name"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	Thickness float64 `json:"thickness"`
	Grain     string  `json:"grain,omitempty"`
	Quantity  int     `json:"quantity"`
}
