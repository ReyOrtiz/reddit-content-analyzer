package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// GetEmbedding Tests
// ============================================================================

func TestClient_GetEmbedding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedEmbedding := []float32{0.1, 0.2, 0.3, 0.4, 0.5}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/embeddings", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var req EmbeddingRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, 1, len(req.Input))
			assert.Equal(t, "test text", req.Input[0])

			response := EmbeddingResponse{
				Data: []struct {
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{
					{
						Embedding: expectedEmbedding,
						Index:     0,
					},
				},
				Model: "text-embedding-mxbai-embed-large-v1",
				Usage: struct {
					PromptTokens int `json:"prompt_tokens"`
					TotalTokens  int `json:"total_tokens"`
				}{
					PromptTokens: 10,
					TotalTokens:  10,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &Client{
			baseURL:        server.URL,
			embeddingModel: "text-embedding-mxbai-embed-large-v1",
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}

		// Act
		result, err := client.GetEmbedding(ctx, "test text")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedEmbedding, result)
		assert.Len(t, result, 5)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client := &Client{
			baseURL:        server.URL,
			embeddingModel: "text-embedding-mxbai-embed-large-v1",
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}

		// Act
		result, err := client.GetEmbedding(ctx, "test text")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "status 500")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := &Client{
			baseURL:        server.URL,
			embeddingModel: "text-embedding-mxbai-embed-large-v1",
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}

		// Act
		result, err := client.GetEmbedding(ctx, "test text")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "decode")
	})

	t.Run("EmptyEmbeddingData", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := EmbeddingResponse{
				Data:  []struct {
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{},
				Model: "text-embedding-mxbai-embed-large-v1",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &Client{
			baseURL:        server.URL,
			embeddingModel: "text-embedding-mxbai-embed-large-v1",
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}

		// Act
		result, err := client.GetEmbedding(ctx, "test text")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no embedding data")
	})

	t.Run("NetworkError", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		client := &Client{
			baseURL:        "http://invalid-url-that-does-not-exist:12345",
			embeddingModel: "text-embedding-mxbai-embed-large-v1",
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}

		// Act
		result, err := client.GetEmbedding(ctx, "test text")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("EmptyText", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req EmbeddingRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "", req.Input[0])

			response := EmbeddingResponse{
				Data: []struct {
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{
					{
						Embedding: []float32{0.0, 0.0, 0.0},
						Index:     0,
					},
				},
				Model: "text-embedding-mxbai-embed-large-v1",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &Client{
			baseURL:        server.URL,
			embeddingModel: "text-embedding-mxbai-embed-large-v1",
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}

		// Act
		result, err := client.GetEmbedding(ctx, "")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// ============================================================================
// Chat Tests
// ============================================================================

func TestClient_Chat(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedResponse := "This is a test response"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/chat/completions", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var req ChatRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, 1, len(req.Messages))
			assert.Equal(t, "user", req.Messages[0].Role)
			assert.Equal(t, "test message", req.Messages[0].Content)

			response := ChatResponse{
				ID:     "chat-123",
				Object: "chat.completion",
				Model:  "openai/gpt-oss-20b",
				Choices: []struct {
					Index        int     `json:"index"`
					Message      Message `json:"message"`
					FinishReason string  `json:"finish_reason"`
				}{
					{
						Index: 0,
						Message: Message{
							Role:    "assistant",
							Content: expectedResponse,
						},
						FinishReason: "stop",
					},
				},
				Usage: struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				}{
					PromptTokens:     10,
					CompletionTokens: 20,
					TotalTokens:      30,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &Client{
			baseURL:   server.URL,
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "user",
				Content: "test message",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
	})

	t.Run("MultipleMessages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedResponse := "Response to multiple messages"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req ChatRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, 2, len(req.Messages))

			response := ChatResponse{
				ID:     "chat-123",
				Object: "chat.completion",
				Model:  "openai/gpt-oss-20b",
				Choices: []struct {
					Index        int     `json:"index"`
					Message      Message `json:"message"`
					FinishReason string  `json:"finish_reason"`
				}{
					{
						Index: 0,
						Message: Message{
							Role:    "assistant",
							Content: expectedResponse,
						},
						FinishReason: "stop",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &Client{
			baseURL:   server.URL,
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "user",
				Content: "first message",
			},
			{
				Role:    "assistant",
				Content: "previous response",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
	})

	t.Run("NoUserMessages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		client := &Client{
			baseURL:   "http://localhost:1234",
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "assistant",
				Content: "only assistant message",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "no user messages")
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		}))
		defer server.Close()

		client := &Client{
			baseURL:   server.URL,
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "user",
				Content: "test message",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "status 400")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := &Client{
			baseURL:   server.URL,
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "user",
				Content: "test message",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "decode")
	})

	t.Run("EmptyChoices", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := ChatResponse{
				ID:     "chat-123",
				Object: "chat.completion",
				Model:  "openai/gpt-oss-20b",
				Choices: []struct {
					Index        int     `json:"index"`
					Message      Message `json:"message"`
					FinishReason string  `json:"finish_reason"`
				}{},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := &Client{
			baseURL:   server.URL,
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "user",
				Content: "test message",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "no choices")
	})

	t.Run("NetworkError", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		client := &Client{
			baseURL:   "http://invalid-url-that-does-not-exist:12345",
			chatModel: "openai/gpt-oss-20b",
			httpClient: &http.Client{},
			logger:    logger.GetLogger(),
		}

		messages := []Message{
			{
				Role:    "user",
				Content: "test message",
			},
		}

		// Act
		result, err := client.Chat(ctx, messages)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
