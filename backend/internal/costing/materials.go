package costing

import (
	"fmt"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

type materialAggregate struct {
	materialID string
	name       string
	areaM2     float64
}

func aggregateMaterialAreas(parts []domain.Part) []materialAggregate {
	byID := map[string]*materialAggregate{}

	for _, part := range parts {
		areaM2 := (part.Width * part.Height) / 1_000_000
		entry, ok := byID[part.Material.ID]
		if !ok {
			byID[part.Material.ID] = &materialAggregate{
				materialID: part.Material.ID,
				name:       part.Material.Name,
				areaM2:     areaM2,
			}
			continue
		}
		entry.areaM2 += areaM2
	}

	result := make([]materialAggregate, 0, len(byID))
	for _, entry := range byID {
		result = append(result, *entry)
	}
	return result
}

func aggregateEdgeBanding(model domain.ManufacturingModel) map[string]float64 {
	lengthsM := map[string]float64{}
	for _, edge := range model.EdgeBanding {
		lengthsM[edge.Material] += edge.Length / 1000
	}
	return lengthsM
}

func buildMaterialLines(parts []domain.Part, opts CostOptions) []domain.MaterialCostLine {
	aggregates := aggregateMaterialAreas(parts)
	lines := make([]domain.MaterialCostLine, 0, len(aggregates))

	for _, agg := range aggregates {
		rate := opts.materialRate(agg.materialID)
		lines = append(lines, domain.MaterialCostLine{
			MaterialID:    agg.materialID,
			Name:          agg.name,
			AreaM2:        round2(agg.areaM2),
			UnitCostPerM2: rate,
			Total:         round2(agg.areaM2 * rate),
		})
	}

	return lines
}

func buildEdgeBandingLines(model domain.ManufacturingModel, opts CostOptions) []domain.MaterialCostLine {
	lines := []domain.MaterialCostLine{}
	for material, lengthM := range aggregateEdgeBanding(model) {
		rate := opts.edgeBandingRate(material)
		lines = append(lines, domain.MaterialCostLine{
			MaterialID:    material,
			Name:          fmt.Sprintf("Tapacanto %s", material),
			AreaM2:        round2(lengthM),
			UnitCostPerM2: rate,
			Total:         round2(lengthM * rate),
		})
	}
	return lines
}
