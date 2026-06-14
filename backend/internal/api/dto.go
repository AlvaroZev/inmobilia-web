package api

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

type ParseAIRequest struct {
	Description string `json:"description"`
	Name        string `json:"name,omitempty"`
}

type SolveRequest struct {
	Room           domain.RoomGeometry           `json:"room"`
	Furniture      domain.FurnitureDefinition    `json:"furniture"`
	Installation   domain.InstallationConstraints `json:"installation"`
}

type ManufacturingRequest struct {
	Resolved domain.ResolvedFurniture `json:"resolved"`
}

type CostRequest struct {
	Model domain.ManufacturingModel `json:"model"`
}

type ExportRequest struct {
	FurnitureName string                    `json:"furnitureName,omitempty"`
	Model         domain.ManufacturingModel `json:"model"`
	Cost          *domain.CostResult        `json:"cost,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
