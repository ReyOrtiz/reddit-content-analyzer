package api

import (
	"fmt"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/config"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/services"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	cfg := config.GetConfig()
	relevanceService := services.NewRelevanceService()
	relevanceHandler := NewRelevanceHandler(relevanceService)

	router := gin.Default()
	router.POST("/v1/reddit/relevance/search", relevanceHandler.GetRelevantPosts)
	router.Run(fmt.Sprintf(":%s", cfg.GetString("api.port")))
}
