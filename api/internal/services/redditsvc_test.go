package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/reddit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// newRedditServiceForTesting creates a redditService with a test client for testing
func newRedditServiceForTesting(baseURL string) *redditService {
	testClient := reddit.NewTestClient(baseURL)
	return &redditService{
		client: *testClient,
		logger: zap.NewNop(),
	}
}

// ============================================================================
// GetPosts Tests
// ============================================================================

func TestRedditService_GetPosts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedResponse := &reddit.RedditResponse{
			Data: reddit.RedditData{
				Children: []reddit.RedditChild{
					{
						Data: reddit.RedditPostData{
							Title:       "Test Post",
							Selftext:    "Test content",
							URL:         "https://reddit.com/r/technology/test",
							Score:       100,
							NumComments: 50,
							CreatedUTC:  float64(time.Now().Unix()),
							Permalink:   "/r/technology/test",
						},
					},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		limit := 5

		// Act
		result, err := service.GetPosts(subreddit, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data.Children, 1)
		assert.Equal(t, "Test Post", result.Data.Children[0].Data.Title)
	})

	t.Run("ClientError", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		limit := 5

		// Act
		result, err := service.GetPosts(subreddit, limit)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		// Arrange
		expectedResponse := &reddit.RedditResponse{
			Data: reddit.RedditData{
				Children: []reddit.RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		limit := 5

		// Act
		result, err := service.GetPosts(subreddit, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Data.Children)
	})
}

// ============================================================================
// SearchPosts Tests
// ============================================================================

func TestRedditService_SearchPosts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedResponse := &reddit.RedditResponse{
			Data: reddit.RedditData{
				Children: []reddit.RedditChild{
					{
						Data: reddit.RedditPostData{
							Title:       "AI Discussion",
							Selftext:    "Discussion about AI",
							URL:         "https://reddit.com/r/technology/ai",
							Score:       200,
							NumComments: 100,
							CreatedUTC:  float64(time.Now().Unix()),
							Permalink:   "/r/technology/ai",
						},
					},
					{
						Data: reddit.RedditPostData{
							Title:       "Machine Learning News",
							Selftext:    "Latest ML updates",
							URL:         "https://reddit.com/r/technology/ml",
							Score:       150,
							NumComments: 75,
							CreatedUTC:  float64(time.Now().Unix()),
							Permalink:   "/r/technology/ml",
						},
					},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		query := "artificial intelligence"
		limit := 5

		// Act
		result, err := service.SearchPosts(subreddit, query, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data.Children, 2)
		assert.Equal(t, "AI Discussion", result.Data.Children[0].Data.Title)
		assert.Equal(t, "Machine Learning News", result.Data.Children[1].Data.Title)
	})

	t.Run("ClientError", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		query := "test query"
		limit := 5

		// Act
		result, err := service.SearchPosts(subreddit, query, limit)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("EmptySearchResults", func(t *testing.T) {
		// Arrange
		expectedResponse := &reddit.RedditResponse{
			Data: reddit.RedditData{
				Children: []reddit.RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		query := "nonexistent topic"
		limit := 5

		// Act
		result, err := service.SearchPosts(subreddit, query, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Data.Children)
	})

	t.Run("SpecialCharactersInQuery", func(t *testing.T) {
		// Arrange
		expectedResponse := &reddit.RedditResponse{
			Data: reddit.RedditData{
				Children: []reddit.RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		service := newRedditServiceForTesting(server.URL)
		subreddit := "technology"
		query := "C++ & Python"
		limit := 5

		// Act
		result, err := service.SearchPosts(subreddit, query, limit)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
