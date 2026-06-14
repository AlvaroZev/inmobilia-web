package export

import (
	"math"
	"time"

	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
)

type BuildOptions struct {
	FurnitureName string
	Cost          *domain.CostResult
}

func BuildBOM(model domain.ManufacturingModel, opts BuildOptions) (domain.BillOfMaterials, error) {
	if result := manufacturing.ValidateManufacturingModel(model); !result.Valid {
		return domain.BillOfMaterials{}, ErrInvalidManufacturingModel
	}

	parts := make([]domain.BOMPartLine, len(model.Parts))
	var totalBoardM2 float64
	for i, part := range model.Parts {
		areaM2 := (part.Width * part.Height) / 1_000_000
		totalBoardM2 += areaM2
		parts[i] = domain.BOMPartLine{
			PartID:         part.ID,
			Name:           part.Name,
			Type:           part.Type,
			VolumeID:       part.VolumeID,
			Width:          part.Width,
			Height:         part.Height,
			Thickness:      part.Thickness,
			MaterialID:     part.Material.ID,
			MaterialName:   part.Material.Name,
			GrainDirection: part.GrainDirection,
			AreaM2:         round3(areaM2),
		}
	}

	hardware := make([]domain.BOMHardwareLine, len(model.Hardware))
	hwQty := 0
	for i, item := range model.Hardware {
		hwQty += item.Quantity
		hardware[i] = domain.BOMHardwareLine{
			HardwareType: item.Type,
			Name:         item.Type,
			Quantity:     item.Quantity,
		}
	}

	if opts.Cost != nil {
		for i, hw := range opts.Cost.Hardware {
			for j := range hardware {
				if hardware[j].HardwareType == hw.HardwareType {
					hardware[j].Name = hw.Name
					hardware[j].UnitCost = hw.UnitCost
					hardware[j].Total = hw.Total
					break
				}
			}
			_ = i
		}
	}

	edgeByMaterial := map[string]float64{}
	for _, edge := range model.EdgeBanding {
		edgeByMaterial[edge.Material] += edge.Length / 1000
	}
	edgeLines := make([]domain.BOMEdgeLine, 0, len(edgeByMaterial))
	var totalEdgeM float64
	for material, length := range edgeByMaterial {
		totalEdgeM += length
		edgeLines = append(edgeLines, domain.BOMEdgeLine{
			Material:     material,
			TotalLengthM: round3(length),
		})
	}

	return domain.BillOfMaterials{
		FurnitureID:   model.FurnitureID,
		FurnitureName: opts.FurnitureName,
		GeneratedAt:   time.Now().UTC(),
		Parts:         parts,
		Hardware:      hardware,
		EdgeBanding:   edgeLines,
		Cost:          opts.Cost,
		Summary: domain.BOMSummary{
			PartCount:     len(parts),
			HardwareCount: hwQty,
			TotalBoardM2:  round3(totalBoardM2),
			TotalEdgeM:    round3(totalEdgeM),
		},
	}, nil
}

func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}
