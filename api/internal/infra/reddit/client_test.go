package reddit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// GetPosts Tests
// ============================================================================

func TestClient_GetPosts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{
					{
						Data: RedditPostData{
							Title:       "Test Post",
							Selftext:    "Test content",
							URL:         "https://reddit.com/r/test/post",
							Score:       100,
							NumComments: 50,
							CreatedUTC:  float64(time.Now().Unix()),
							Permalink:   "/r/test/post",
						},
					},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/r/technology/.json?limit=5", r.URL.Path+"?"+r.URL.RawQuery)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "reddit-content-analyzer/1.0", r.Header.Get("User-Agent"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.GetPosts("technology", 5)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data.Children, 1)
		assert.Equal(t, "Test Post", result.Data.Children[0].Data.Title)
	})

	t.Run("LimitDefaulting", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "limit=25")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.GetPosts("technology", 0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("LimitCapping", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "limit=100")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.GetPosts("technology", 200)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.GetPosts("technology", 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "status 500")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.GetPosts("technology", 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "decode")
	})

	t.Run("NetworkError", func(t *testing.T) {
		// Arrange
		client := NewTestClient("http://invalid-url-that-does-not-exist:12345")

		// Act
		result, err := client.GetPosts("technology", 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.GetPosts("nonexistent", 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "status 404")
	})
}

// ============================================================================
// SearchPosts Tests
// ============================================================================

func TestClient_SearchPosts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{
					{
						Data: RedditPostData{
							Title:       "Search Result",
							Selftext:    "Search content",
							URL:         "https://reddit.com/r/test/search",
							Score:       50,
							NumComments: 25,
							CreatedUTC:  float64(time.Now().Unix()),
							Permalink:   "/r/test/search",
						},
					},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/r/technology/search.json")
			assert.Contains(t, r.URL.RawQuery, "q=artificial+intelligence")
			assert.Contains(t, r.URL.RawQuery, "restrict_sr=true")
			assert.Contains(t, r.URL.RawQuery, "limit=5")
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "reddit-content-analyzer/1.0", r.Header.Get("User-Agent"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "artificial intelligence", 5)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data.Children, 1)
		assert.Equal(t, "Search Result", result.Data.Children[0].Data.Title)
	})

	t.Run("LimitDefaulting", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "limit=25")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "test", 0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("LimitCapping", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "limit=100")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "test", 150)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("QueryEncoding", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check that special characters are properly encoded
			assert.Contains(t, r.URL.RawQuery, "q=")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "C++ & Python", 5)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "test", 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "status 400")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("not json"))
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "test", 5)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "decode")
	})

	t.Run("EmptyResults", func(t *testing.T) {
		// Arrange
		expectedResponse := &RedditResponse{
			Data: RedditData{
				Children: []RedditChild{},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTestClient(server.URL)

		// Act
		result, err := client.SearchPosts("technology", "nonexistent", 5)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Data.Children)
	})
}
