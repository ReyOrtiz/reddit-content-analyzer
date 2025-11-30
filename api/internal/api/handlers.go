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

// GetRelevantPosts godoc
// @Summary      Search for relevant Reddit posts
// @Description  Searches Reddit posts based on a topic and returns posts that are relevant according to the specified criteria
// @Tags         reddit
// @Accept       json
// @Produce      json
// @Param        request  body      contracts.RelevanceRequestDto  true  "Search request parameters"
// @Success      200      {object}  contracts.RelevanceResponseDto  "Successful response with relevant posts"
// @Failure      400      {object}  map[string]string              "Bad request - invalid input parameters"
// @Failure      500      {object}  map[string]string              "Internal server error"
// @Router       /v1/reddit/relevance/search [post]
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
