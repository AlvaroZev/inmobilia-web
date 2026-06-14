package api

import (
	"github.com/gin-gonic/gin"
	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
)

func NewRouter(parser ai.Parser, info ServiceInfo) *gin.Engine {
	router := gin.New()
	_ = router.SetTrustedProxies(trustedProxies())
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	handler := NewHandler(parser, info)

	router.GET("/health", handler.Health)
	router.POST("/ai/parse", handler.ParseAI)
	router.POST("/solver", handler.Solve)
	router.POST("/manufacturing", handler.CompileManufacturing)
	router.POST("/cost", handler.CalculateCost)
	router.POST("/export/bom", handler.ExportBOM)
	router.POST("/export/cut-plans", handler.ExportCutPlan)
	router.POST("/export/pdf", handler.ExportPDF)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigin := corsAllowedOrigin()
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
