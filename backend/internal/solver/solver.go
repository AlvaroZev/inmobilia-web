package solver

import (
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
	"github.com/inmobilia/inmobilia-web/backend/internal/resolvedfurniture"
	"github.com/inmobilia/inmobilia-web/backend/internal/roomgeometry"
	"github.com/inmobilia/inmobilia-web/backend/internal/volumetree"
)

// SolveConstraints resolves abstract furniture definition into real dimensions.
func SolveConstraints(
	room domain.RoomGeometry,
	furniture domain.FurnitureDefinition,
	installation domain.InstallationConstraints,
) (domain.ResolvedFurniture, error) {
	if result := roomgeometry.ValidateRoomGeometry(room); !result.Valid {
		return domain.ResolvedFurniture{}, ErrInvalidRoom
	}

	if result := volumetree.ValidateFurnitureDefinition(furniture); !result.Valid {
		return domain.ResolvedFurniture{}, ErrInvalidFurniture
	}

	installBox, err := computeInstallSpace(room, furniture, installation)
	if err != nil {
		return domain.ResolvedFurniture{}, err
	}

	root, err := resolveVolumeTree(furniture.Root, installBox, nil)
	if err != nil {
		return domain.ResolvedFurniture{}, err
	}

	resolved := domain.ResolvedFurniture{
		ID:   furniture.ID,
		Name: furniture.Name,
		Root: root,
	}

	if result := resolvedfurniture.ValidateResolvedFurniture(resolved); !result.Valid {
		return domain.ResolvedFurniture{}, resolvedfurniture.ErrInvalidResolvedFurniture
	}

	return resolved, nil
}
