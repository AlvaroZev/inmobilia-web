package costing

import (
	"math"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func sumMaterialTotals(lines []domain.MaterialCostLine) float64 {
	var total float64
	for _, line := range lines {
		total += line.Total
	}
	return total
}

func sumHardwareTotals(lines []domain.HardwareCostLine) float64 {
	var total float64
	for _, line := range lines {
		total += line.Total
	}
	return total
}
