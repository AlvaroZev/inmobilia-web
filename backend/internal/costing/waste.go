package costing

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

func buildWasteCost(materialLines []domain.MaterialCostLine, opts CostOptions) domain.WasteCost {
	var boardAreaM2 float64
	var boardCost float64

	for _, line := range materialLines {
		if line.MaterialID == "pvc-white-1mm" {
			continue
		}
		boardAreaM2 += line.AreaM2
		boardCost += line.Total
	}

	return domain.WasteCost{
		AreaM2:     round2(boardAreaM2 * opts.WastePercentage),
		Percentage: round2(opts.WastePercentage * 100),
		Total:      round2(boardCost * opts.WastePercentage),
	}
}
