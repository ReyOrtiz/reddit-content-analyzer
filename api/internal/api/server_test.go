package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/contracts"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	mock_services "github.com/ReyOrtiz/reddit-content-analyzer/mocks/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Server Setup Tests
// ============================================================================

func TestStartServer(t *testing.T) {
	// Note: This is a basic test structure. Full server testing would require
	// more setup and potentially integration testing.
	// This demonstrates the test pattern for server initialization.

	t.Run("ServerInitialization", func(t *testing.T) {
		// This test would require refactoring StartServer to accept dependencies
		// For now, it demonstrates the structure
		_ = StartServer
	})
}

// ============================================================================
// Handler Tests
// ============================================================================

func TestRelevanceHandler_GetRelevantPosts(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRelevanceService := mock_services.NewMockRelevanceService(t)
		handler := &RelevanceHandler{
			logger:           logger.GetLogger(),
			relevanceService: mockRelevanceService,
		}

		request := contracts.RelevanceRequestDto{
			Topic:              "artificial intelligence",
			Subreddits:         []string{"technology"},
			RelevanceThreshold: 0.7,
			Limit:              5,
			SearchMethod:       contracts.SearchMethodSearch,
		}

		expectedResponse := contracts.RelevanceResponseDto{
			Posts: []contracts.SubRedditPostDto{
				{
					SubredditName:    "technology",
					Title:            "AI Post",
					Content:          "Content about AI",
					Url:              "https://reddit.com/r/technology/ai",
					Score:            100,
					NumComments:      50,
					RelevanceScore:   0.85,
					IsRelevant:       true,
					RelevanceSummary: "This post is highly relevant",
				},
			},
		}

		mockRelevanceService.EXPECT().
			GetRelevantPosts(mock.Anything, request).
			Return(expectedResponse, nil)

		// Create request body
		requestBody, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/v1/reddit/relevance/search", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler.GetRelevantPosts(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response contracts.RelevanceResponseDto
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Posts, 1)
		assert.Equal(t, "technology", response.Posts[0].SubredditName)
		assert.Equal(t, "AI Post", response.Posts[0].Title)
		assert.True(t, response.Posts[0].IsRelevant)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockRelevanceService := mock_services.NewMockRelevanceService(t)
		handler := &RelevanceHandler{
			logger:           logger.GetLogger(),
			relevanceService: mockRelevanceService,
		}

		// Create invalid request body
		req, _ := http.NewRequest("POST", "/v1/reddit/relevance/search", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler.GetRelevantPosts(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		// Arrange
		mockRelevanceService := mock_services.NewMockRelevanceService(t)
		handler := &RelevanceHandler{
			logger:           logger.GetLogger(),
			relevanceService: mockRelevanceService,
		}

		// Create request with missing fields
		invalidRequest := map[string]interface{}{
			"topic": "test",
			// Missing subreddits, threshold, etc.
		}
		requestBody, _ := json.Marshal(invalidRequest)
		req, _ := http.NewRequest("POST", "/v1/reddit/relevance/search", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler.GetRelevantPosts(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Arrange
		mockRelevanceService := mock_services.NewMockRelevanceService(t)
		handler := &RelevanceHandler{
			logger:           logger.GetLogger(),
			relevanceService: mockRelevanceService,
		}

		request := contracts.RelevanceRequestDto{
			Topic:              "test topic",
			Subreddits:         []string{"test"},
			RelevanceThreshold: 0.7,
			Limit:              5,
			SearchMethod:       contracts.SearchMethodSearch,
		}

		mockRelevanceService.EXPECT().
			GetRelevantPosts(mock.Anything, request).
			Return(contracts.RelevanceResponseDto{}, assert.AnError)

		requestBody, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "/v1/reddit/relevance/search", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler.GetRelevantPosts(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("EmptyRequest", func(t *testing.T) {
		// Arrange
		mockRelevanceService := mock_services.NewMockRelevanceService(t)
		handler := &RelevanceHandler{
			logger:           logger.GetLogger(),
			relevanceService: mockRelevanceService,
		}

		req, _ := http.NewRequest("POST", "/v1/reddit/relevance/search", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler.GetRelevantPosts(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNewRelevanceHandler(t *testing.T) {
	t.Run("CreatesHandler", func(t *testing.T) {
		// Arrange
		mockRelevanceService := mock_services.NewMockRelevanceService(t)

		// Act
		handler := NewRelevanceHandler(mockRelevanceService)

		// Assert
		assert.NotNil(t, handler)
		assert.Equal(t, mockRelevanceService, handler.relevanceService)
	})
}

