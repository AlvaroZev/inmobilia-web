package domain

// CostResult is the output of the cost engine.

type CostResult struct {
	FurnitureID string             `json:"furnitureId"`
	Currency    string             `json:"currency"`
	Materials   []MaterialCostLine `json:"materials"`
	Hardware    []HardwareCostLine `json:"hardware"`
	Labor       LaborCost          `json:"labor"`
	Waste       WasteCost          `json:"waste"`
	Subtotal    float64            `json:"subtotal"`
	Total       float64            `json:"total"`
}

type MaterialCostLine struct {
	MaterialID      string  `json:"materialId"`
	Name            string  `json:"name"`
	AreaM2          float64 `json:"areaM2"`
	UnitCostPerM2   float64 `json:"unitCostPerM2"`
	Total           float64 `json:"total"`
}

type HardwareCostLine struct {
	HardwareType string  `json:"hardwareType"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	UnitCost     float64 `json:"unitCost"`
	Total        float64 `json:"total"`
}

type LaborCost struct {
	Hours       float64 `json:"hours"`
	RatePerHour float64 `json:"ratePerHour"`
	Total       float64 `json:"total"`
}

type WasteCost struct {
	AreaM2     float64 `json:"areaM2"`
	Percentage float64 `json:"percentage"`
	Total      float64 `json:"total"`
}
