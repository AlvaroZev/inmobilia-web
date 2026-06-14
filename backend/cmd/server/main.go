package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inmobilia/inmobilia-web/backend/internal/api"
	"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	parser, providerLabel := ai.ResolveParserFromEnv()
	log.Printf("AI provider: %s", providerLabel)

	router := api.NewRouter(parser, api.ServiceInfo{
		AIProvider: providerLabel,
	})

	addr := fmt.Sprintf(":%s", port)
	if err := router.Run(addr); err != nil {
		panic(err)
	}
}
