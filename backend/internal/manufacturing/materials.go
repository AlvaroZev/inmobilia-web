package manufacturing

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

const (
	defaultBoardThickness = 18
	defaultBackThickness  = 3
	defaultEdgeBanding    = "pvc-white-1mm"
)

func resolveMaterial(materialID string) domain.Material {
	switch materialID {
	case "melamine-white-18", "melamine-white":
		return domain.Material{
			ID:        "melamine-white-18",
			Name:      "Melamina Blanca 18mm",
			Type:      "melamine",
			Thickness: defaultBoardThickness,
			Color:     "white",
		}
	case "nordex":
		return domain.Material{
			ID:        "nordex",
			Name:      "Nordex",
			Type:      "nordex",
			Thickness: defaultBackThickness,
			Color:     "natural",
		}
	default:
		id := materialID
		if id == "" {
			id = "melamine-default-18"
		}
		return domain.Material{
			ID:        id,
			Name:      "Melamina 18mm",
			Type:      "melamine",
			Thickness: defaultBoardThickness,
		}
	}
}

func resolveEdgeBanding(materialID string) string {
	if materialID == "" {
		return defaultEdgeBanding
	}
	return defaultEdgeBanding
}

func backMaterial(board domain.Material) domain.Material {
	return domain.Material{
		ID:        board.ID + "-back",
		Name:      board.Name + " (trasera)",
		Type:      "back_panel",
		Thickness: defaultBackThickness,
		Color:     board.Color,
	}
}
