package export

import (
	"sort"
	"time"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
)

func BuildCutPlan(model domain.ManufacturingModel, opts BuildOptions) (domain.CutPlan, error) {
	if result := manufacturing.ValidateManufacturingModel(model); !result.Valid {
		return domain.CutPlan{}, ErrInvalidManufacturingModel
	}

	type sheetKey struct {
		materialID string
		thickness  float64
	}

	sheets := map[sheetKey]*domain.CutSheet{}

	for _, part := range model.Parts {
		key := sheetKey{materialID: part.Material.ID, thickness: part.Material.Thickness}
		sheet, ok := sheets[key]
		if !ok {
			sheet = &domain.CutSheet{
				MaterialID:   part.Material.ID,
				MaterialName: part.Material.Name,
				Thickness:    part.Material.Thickness,
				Parts:        []domain.CutPartLine{},
			}
			sheets[key] = sheet
		}

		areaM2 := (part.Width * part.Height) / 1_000_000
		sheet.TotalAreaM2 += areaM2
		sheet.Parts = append(sheet.Parts, domain.CutPartLine{
			PartID:    part.ID,
			Name:      part.Name,
			Width:     part.Width,
			Height:    part.Height,
			Thickness: part.Thickness,
			Grain:     part.GrainDirection,
			Quantity:  1,
		})
	}

	resultSheets := make([]domain.CutSheet, 0, len(sheets))
	for _, sheet := range sheets {
		sort.Slice(sheet.Parts, func(i, j int) bool {
			areaI := sheet.Parts[i].Width * sheet.Parts[i].Height
			areaJ := sheet.Parts[j].Width * sheet.Parts[j].Height
			return areaI > areaJ
		})
		sheet.TotalAreaM2 = round3(sheet.TotalAreaM2)
		resultSheets = append(resultSheets, *sheet)
	}

	sort.Slice(resultSheets, func(i, j int) bool {
		return resultSheets[i].MaterialName < resultSheets[j].MaterialName
	})

	return domain.CutPlan{
		FurnitureID:   model.FurnitureID,
		FurnitureName: opts.FurnitureName,
		GeneratedAt:   time.Now().UTC(),
		Sheets:        resultSheets,
	}, nil
}
