package costing

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

var hardwareNames = map[string]string{
	"hinge":         "Bisagra",
	"drawer_runner": "Corredera",
	"hanger_rod":    "Barral",
	"rod_bracket":   "Soporte barral",
}

func buildHardwareLines(hardware []domain.Hardware, opts CostOptions) []domain.HardwareCostLine {
	byType := map[string]*domain.HardwareCostLine{}

	for _, item := range hardware {
		line, ok := byType[item.Type]
		if !ok {
			name := hardwareNames[item.Type]
			if name == "" {
				name = item.Type
			}
			byType[item.Type] = &domain.HardwareCostLine{
				HardwareType: item.Type,
				Name:         name,
				Quantity:     0,
				UnitCost:     opts.hardwareRate(item.Type),
			}
			line = byType[item.Type]
		}
		line.Quantity += item.Quantity
	}

	lines := make([]domain.HardwareCostLine, 0, len(byType))
	for _, line := range byType {
		line.Total = round2(float64(line.Quantity) * line.UnitCost)
		lines = append(lines, *line)
	}
	return lines
}
