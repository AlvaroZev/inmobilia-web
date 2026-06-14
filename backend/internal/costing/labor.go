package costing

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

func estimateLaborHours(model domain.ManufacturingModel) float64 {
	hours := 1.5
	hours += float64(len(model.Parts)) * 0.12

	for _, part := range model.Parts {
		if part.Type == string(domain.PartDoor) {
			hours += 0.4
		}
	}

	for _, hw := range model.Hardware {
		switch hw.Type {
		case "drawer_runner":
			hours += float64(hw.Quantity) * 0.35
		case "hinge":
			hours += float64(hw.Quantity) * 0.08
		}
	}

	return round2(hours)
}

func buildLaborCost(model domain.ManufacturingModel, opts CostOptions) domain.LaborCost {
	hours := estimateLaborHours(model)
	return domain.LaborCost{
		Hours:       hours,
		RatePerHour: opts.LaborRatePerHour,
		Total:       round2(hours * opts.LaborRatePerHour),
	}
}
