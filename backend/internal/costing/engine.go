package costing

import (
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
)

// CalculateCost computes materials, hardware, labor, waste, and total cost.
func CalculateCost(model domain.ManufacturingModel) (domain.CostResult, error) {
	return CalculateCostWithOptions(model, DefaultCostOptions())
}

// CalculateCostWithOptions computes cost using custom unit rates.
func CalculateCostWithOptions(model domain.ManufacturingModel, opts CostOptions) (domain.CostResult, error) {
	if result := manufacturing.ValidateManufacturingModel(model); !result.Valid {
		return domain.CostResult{}, ErrInvalidManufacturingModel
	}

	boardLines := buildMaterialLines(model.Parts, opts)
	edgeLines := buildEdgeBandingLines(model, opts)
	materialLines := append(boardLines, edgeLines...)

	hardwareLines := buildHardwareLines(model.Hardware, opts)
	labor := buildLaborCost(model, opts)
	waste := buildWasteCost(boardLines, opts)

	subtotal := round2(
		sumMaterialTotals(materialLines) +
			sumHardwareTotals(hardwareLines) +
			labor.Total,
	)

	cost := domain.CostResult{
		FurnitureID: model.FurnitureID,
		Currency:    opts.Currency,
		Materials:   materialLines,
		Hardware:    hardwareLines,
		Labor:       labor,
		Waste:       waste,
		Subtotal:    subtotal,
		Total:       round2(subtotal + waste.Total),
	}

	if result := ValidateCostResult(cost); !result.Valid {
		return domain.CostResult{}, ErrInvalidCostResult
	}

	return cost, nil
}
