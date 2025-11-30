package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/config"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

var (
	client *Client
	once   sync.Once
)

// ClientInterface defines the interface for LLM client operations
type ClientInterface interface {
	GetEmbedding(ctx context.Context, text string) ([]float32, error)
	Chat(ctx context.Context, messages []Message) (string, error)
}

// Client represents an LLM client using Genkit Go
type Client struct {
	genkit         *genkit.Genkit
	baseURL        string
	embeddingModel string
	chatModel      string
	httpClient     *http.Client
	logger         *zap.Logger
}

// EmbeddingRequest represents a request for embeddings
type EmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// EmbeddingResponse represents the response from the embedding API
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatRequest represents a request for chat completion
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the response from the chat API
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GetClient returns the singleton LLM client instance, initializing it on first call
func GetClient() *Client {
	once.Do(func() {
		cfg := config.GetConfig()
		baseURL := cfg.GetString("llm.base_url")
		embeddingModel := cfg.GetString("llm.embedding_model")
		chatModel := cfg.GetString("llm.summarization_model")

		if baseURL == "" {
			baseURL = "http://127.0.0.1:1234/v1"
		}
		if embeddingModel == "" {
			embeddingModel = "text-embedding-mxbai-embed-large-v1"
		}
		if chatModel == "" {
			chatModel = "openai/gpt-oss-20b"
		}

		ctx := context.Background()
		g := genkit.Init(ctx)

		client = &Client{
			genkit:         g,
			baseURL:        baseURL,
			embeddingModel: embeddingModel,
			chatModel:      chatModel,
			httpClient:     &http.Client{},
			logger:         logger.GetLogger(),
		}
	})
	return client
}

// GetEmbedding generates embeddings for the given text using the configured embedding model
func (c *Client) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	c.logger.Info("Generating embedding", zap.String("text", text), zap.String("model", c.embeddingModel))

	// Use OpenAI-compatible API for embeddings
	url := fmt.Sprintf("%s/embeddings", c.baseURL)
	req := EmbeddingRequest{
		Input: []string{text},
		Model: c.embeddingModel,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Error marshaling embedding request", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Error creating embedding request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Error calling embedding API", zap.Error(err))
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Embedding API returned error", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		c.logger.Error("Error decoding embedding response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	c.logger.Info("Embedding generated successfully", zap.Int("dimension", len(embeddingResp.Data[0].Embedding)))
	return embeddingResp.Data[0].Embedding, nil
}

// Chat sends a chat message and returns the model's response
func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	c.logger.Info("Sending chat message", zap.String("model", c.chatModel), zap.Int("message_count", len(messages)))

	// Use Genkit's Generate function for chat
	// Convert messages to Genkit's format
	var promptParts []*ai.Part
	for _, msg := range messages {
		if msg.Role == "user" {
			promptParts = append(promptParts, ai.NewTextPart(msg.Content))
		}
	}

	if len(promptParts) == 0 {
		return "", fmt.Errorf("no user messages found")
	}

	// Use Genkit's Generate with the configured model
	// For OpenAI-compatible APIs, we'll use HTTP directly
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	req := ChatRequest{
		Model:    c.chatModel,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Error marshaling chat request", zap.Error(err))
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Error creating chat request", zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Error calling chat API", zap.Error(err))
		return "", fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Chat API returned error", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		c.logger.Error("Error decoding chat response", zap.Error(err))
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	responseText := chatResp.Choices[0].Message.Content
	c.logger.Info("Chat response received", zap.String("response", responseText))
	return responseText, nil
}
