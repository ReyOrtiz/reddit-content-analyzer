package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/contracts"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/services"
)

type RelevanceHandler struct {
	logger           *zap.Logger
	relevanceService services.RelevanceService
}

func NewRelevanceHandler(relevanceService services.RelevanceService) *RelevanceHandler {
	return &RelevanceHandler{
		logger:           logger.GetLogger(),
		relevanceService: relevanceService,
	}
}

func (h *RelevanceHandler) GetRelevantPosts(c *gin.Context) {
	h.logger.Info("Searching Reddit posts", zap.Any("request", c.Request.Body))

	var request contracts.RelevanceRequestDto
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.relevanceService.GetRelevantPosts(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Error searching Reddit posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
