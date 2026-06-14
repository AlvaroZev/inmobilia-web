package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inmobilia/inmobilia-web/backend/internal/costing"
	"github.com/inmobilia/inmobilia-web/backend/internal/export"
	"github.com/inmobilia/inmobilia-web/backend/internal/manufacturing"
	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
	"github.com/inmobilia/inmobilia-web/backend/internal/solver"
)

type Handler struct {
	parser ai.Parser
	info   ServiceInfo
}

func NewHandler(parser ai.Parser, info ServiceInfo) *Handler {
	return &Handler{parser: parser, info: info}
}

func (h *Handler) ParseAI(c *gin.Context) {
	var req ParseAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	furniture, err := h.parser.ParseFurniture(c.Request.Context(), req.Description, req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ai.ErrOpenAIRequest) || errors.Is(err, ai.ErrOpenAINotConfigured) {
			status = http.StatusBadGateway
		}
		c.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, furniture)
}

func (h *Handler) Solve(c *gin.Context) {
	var req SolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	resolved, err := solver.SolveConstraints(req.Room, req.Furniture, req.Installation)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resolved)
}

func (h *Handler) CompileManufacturing(c *gin.Context) {
	var req ManufacturingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	model, err := manufacturing.CompileManufacturing(req.Resolved)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) CalculateCost(c *gin.Context) {
	var req CostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	result, err := costing.CalculateCost(req.Model)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) ExportBOM(c *gin.Context) {
	req, opts, ok := bindExportRequest(c)
	if !ok {
		return
	}

	bom, err := export.BuildBOM(req.Model, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bom)
}

func (h *Handler) ExportCutPlan(c *gin.Context) {
	req, opts, ok := bindExportRequest(c)
	if !ok {
		return
	}

	plan, err := export.BuildCutPlan(req.Model, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) ExportPDF(c *gin.Context) {
	req, opts, ok := bindExportRequest(c)
	if !ok {
		return
	}

	bom, err := export.BuildBOM(req.Model, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	plan, err := export.BuildCutPlan(req.Model, opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	pdf, err := export.GeneratePDF(bom, plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", `attachment; filename="inmobilia-`+req.Model.FurnitureID+`.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdf)
}

func bindExportRequest(c *gin.Context) (ExportRequest, export.BuildOptions, bool) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return req, export.BuildOptions{}, false
	}
	return req, export.BuildOptions{
		FurnitureName: req.FurnitureName,
		Cost:          req.Cost,
	}, true
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"service":     "inmobilia-api",
		"ai_provider": h.info.AIProvider,
	})
}
