package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/contracts"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/llm"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/reddit"
	mock_llm "github.com/ReyOrtiz/reddit-content-analyzer/mocks/llm"
	mock_services "github.com/ReyOrtiz/reddit-content-analyzer/mocks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// newRelevanceServiceForTesting creates a relevanceService with injected dependencies for testing
func newRelevanceServiceForTesting(llmClient llm.ClientInterface, redditService RedditService) *relevanceService {
	return &relevanceService{
		logger:        zap.NewNop(),
		llmClient:     llmClient,
		redditService: redditService,
	}
}

// ============================================================================
// GetRelevantPosts Tests
// ============================================================================

func TestRelevanceService_GetRelevantPosts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("SearchMethod", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			topic := "artificial intelligence"
			subreddit := "technology"
			limit := 5
			relevanceThreshold := 0.7

			request := contracts.RelevanceRequestDto{
				Topic:              topic,
				Subreddits:         []string{subreddit},
				RelevanceThreshold: relevanceThreshold,
				Limit:              limit,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			topicEmbedding := []float32{0.1, 0.2, 0.3, 0.4, 0.5}
			post1Embedding := []float32{0.12, 0.22, 0.32, 0.42, 0.52} // High similarity (>0.7)
			post2Embedding := []float32{1.0, 0.0, 0.0, 0.0, 0.0}      // Low similarity (<0.7)

			redditResponse := &reddit.RedditResponse{
				Data: reddit.RedditData{
					Children: []reddit.RedditChild{
						{
							Data: reddit.RedditPostData{
								Title:       "AI in Healthcare",
								Selftext:    "Discussion about AI applications in healthcare",
								URL:         "https://reddit.com/r/technology/ai-healthcare",
								Score:       100,
								NumComments: 50,
								CreatedUTC:  float64(time.Now().Unix()),
								Permalink:   "/r/technology/ai-healthcare",
							},
						},
						{
							Data: reddit.RedditPostData{
								Title:       "Random Post",
								Selftext:    "This is unrelated content",
								URL:         "https://reddit.com/r/technology/random",
								Score:       10,
								NumComments: 5,
								CreatedUTC:  float64(time.Now().Unix()),
								Permalink:   "/r/technology/random",
							},
						},
					},
				},
			}

			mockLLMClient.EXPECT().GetEmbedding(ctx, topic).Return(topicEmbedding, nil)
			mockRedditService.EXPECT().SearchPosts(subreddit, topic, limit).Return(redditResponse, nil)
			mockLLMClient.EXPECT().GetEmbedding(ctx, "AI in Healthcare. Discussion about AI applications in healthcare").
				Return(post1Embedding, nil)
			mockLLMClient.EXPECT().GetEmbedding(ctx, "Random Post. This is unrelated content").
				Return(post2Embedding, nil)
			mockLLMClient.EXPECT().Chat(ctx, mock.MatchedBy(func(messages []llm.Message) bool {
				return len(messages) == 1 && messages[0].Role == "user"
			})).Return("This post is highly relevant to artificial intelligence", nil).Times(2)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.NoError(t, err)
			assert.Len(t, result.Posts, 2)

			post1 := result.Posts[0]
			assert.Equal(t, subreddit, post1.SubredditName)
			assert.Equal(t, "AI in Healthcare", post1.Title)
			assert.True(t, post1.IsRelevant)
			assert.Greater(t, post1.RelevanceScore, relevanceThreshold)
			assert.NotEmpty(t, post1.RelevanceSummary)

			post2 := result.Posts[1]
			assert.Equal(t, subreddit, post2.SubredditName)
			assert.Equal(t, "Random Post", post2.Title)
			assert.False(t, post2.IsRelevant)
			assert.Less(t, post2.RelevanceScore, relevanceThreshold)
			assert.NotEmpty(t, post2.RelevanceSummary)
		})

		t.Run("LatestMethod", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			topic := "machine learning"
			subreddit := "MachineLearning"
			limit := 3
			relevanceThreshold := 0.6

			request := contracts.RelevanceRequestDto{
				Topic:              topic,
				Subreddits:         []string{subreddit},
				RelevanceThreshold: relevanceThreshold,
				Limit:              limit,
				SearchMethod:       contracts.SearchMethodLatest,
			}

			topicEmbedding := []float32{0.2, 0.3, 0.4, 0.5, 0.6}
			postEmbedding := []float32{0.25, 0.35, 0.45, 0.55, 0.65}

			redditResponse := &reddit.RedditResponse{
				Data: reddit.RedditData{
					Children: []reddit.RedditChild{
						{
							Data: reddit.RedditPostData{
								Title:       "New ML Paper",
								Selftext:    "Latest research in machine learning",
								URL:         "https://reddit.com/r/MachineLearning/new-paper",
								Score:       200,
								NumComments: 100,
								CreatedUTC:  float64(time.Now().Unix()),
								Permalink:   "/r/MachineLearning/new-paper",
							},
						},
					},
				},
			}

			mockLLMClient.EXPECT().GetEmbedding(ctx, topic).Return(topicEmbedding, nil)
			mockRedditService.EXPECT().GetPosts(subreddit, limit).Return(redditResponse, nil)
			mockLLMClient.EXPECT().GetEmbedding(ctx, "New ML Paper. Latest research in machine learning").
				Return(postEmbedding, nil)
			mockLLMClient.EXPECT().Chat(ctx, mock.Anything).Return("This post discusses machine learning research", nil)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.NoError(t, err)
			assert.Len(t, result.Posts, 1)
			assert.Equal(t, subreddit, result.Posts[0].SubredditName)
			assert.Equal(t, "New ML Paper", result.Posts[0].Title)
		})

		t.Run("MultipleSubreddits", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			topic := "programming"
			subreddits := []string{"programming", "golang"}
			limit := 2
			relevanceThreshold := 0.5

			request := contracts.RelevanceRequestDto{
				Topic:              topic,
				Subreddits:         subreddits,
				RelevanceThreshold: relevanceThreshold,
				Limit:              limit,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			topicEmbedding := []float32{0.1, 0.2, 0.3}
			postEmbedding := []float32{0.15, 0.25, 0.35}

			for _, subreddit := range subreddits {
				redditResponse := &reddit.RedditResponse{
					Data: reddit.RedditData{
						Children: []reddit.RedditChild{
							{
								Data: reddit.RedditPostData{
									Title:       "Post in " + subreddit,
									Selftext:    "Content about " + topic,
									URL:         "https://reddit.com/r/" + subreddit + "/post",
									Score:       50,
									NumComments: 25,
									CreatedUTC:  float64(time.Now().Unix()),
									Permalink:   "/r/" + subreddit + "/post",
								},
							},
						},
					},
				}

				mockRedditService.EXPECT().SearchPosts(subreddit, topic, limit).Return(redditResponse, nil)
				mockLLMClient.EXPECT().GetEmbedding(ctx, mock.MatchedBy(func(text string) bool {
					return len(text) > 0
				})).Return(postEmbedding, nil)
				mockLLMClient.EXPECT().Chat(ctx, mock.Anything).Return("Relevant post about programming", nil)
			}

			mockLLMClient.EXPECT().GetEmbedding(ctx, topic).Return(topicEmbedding, nil)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.NoError(t, err)
			assert.Len(t, result.Posts, 2)
			assert.Equal(t, subreddits[0], result.Posts[0].SubredditName)
			assert.Equal(t, subreddits[1], result.Posts[1].SubredditName)
		})

		t.Run("EmptySubreddits", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			request := contracts.RelevanceRequestDto{
				Topic:              "test topic",
				Subreddits:         []string{},
				RelevanceThreshold: 0.7,
				Limit:              5,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			topicEmbedding := []float32{0.1, 0.2, 0.3}
			mockLLMClient.EXPECT().GetEmbedding(ctx, "test topic").Return(topicEmbedding, nil)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.NoError(t, err)
			assert.Empty(t, result.Posts)
		})
	})

	t.Run("Failure", func(t *testing.T) {
		t.Run("LLMEmbeddingError", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			request := contracts.RelevanceRequestDto{
				Topic:              "test topic",
				Subreddits:         []string{"test"},
				RelevanceThreshold: 0.7,
				Limit:              5,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			expectedError := errors.New("LLM service unavailable")
			mockLLMClient.EXPECT().GetEmbedding(ctx, "test topic").Return(nil, expectedError)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "error getting topic embedding")
			assert.Empty(t, result.Posts)
		})

		t.Run("RedditServiceError", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			request := contracts.RelevanceRequestDto{
				Topic:              "test topic",
				Subreddits:         []string{"test"},
				RelevanceThreshold: 0.7,
				Limit:              5,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			topicEmbedding := []float32{0.1, 0.2, 0.3}
			expectedError := errors.New("Reddit API error")

			mockLLMClient.EXPECT().GetEmbedding(ctx, "test topic").Return(topicEmbedding, nil)
			mockRedditService.EXPECT().SearchPosts("test", "test topic", 5).Return(nil, expectedError)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "error getting subreddit posts")
			assert.Empty(t, result.Posts)
		})

		t.Run("PostEmbeddingError", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			request := contracts.RelevanceRequestDto{
				Topic:              "test topic",
				Subreddits:         []string{"test"},
				RelevanceThreshold: 0.7,
				Limit:              5,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			topicEmbedding := []float32{0.1, 0.2, 0.3}
			expectedError := errors.New("embedding generation failed")

			redditResponse := &reddit.RedditResponse{
				Data: reddit.RedditData{
					Children: []reddit.RedditChild{
						{
							Data: reddit.RedditPostData{
								Title:       "Test Post",
								Selftext:    "Test content",
								URL:         "https://reddit.com/r/test/post",
								Score:       10,
								NumComments: 5,
								CreatedUTC:  float64(time.Now().Unix()),
								Permalink:   "/r/test/post",
							},
						},
					},
				},
			}

			mockLLMClient.EXPECT().GetEmbedding(ctx, "test topic").Return(topicEmbedding, nil)
			mockRedditService.EXPECT().SearchPosts("test", "test topic", 5).Return(redditResponse, nil)
			mockLLMClient.EXPECT().GetEmbedding(ctx, "Test Post. Test content").Return(nil, expectedError)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "error getting relevance score")
			assert.Empty(t, result.Posts)
		})

		t.Run("SummaryError", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockLLMClient := mock_llm.NewMockClientInterface(t)
			mockRedditService := mock_services.NewMockRedditService(t)
			service := newRelevanceServiceForTesting(mockLLMClient, mockRedditService)

			request := contracts.RelevanceRequestDto{
				Topic:              "test topic",
				Subreddits:         []string{"test"},
				RelevanceThreshold: 0.7,
				Limit:              5,
				SearchMethod:       contracts.SearchMethodSearch,
			}

			topicEmbedding := []float32{0.1, 0.2, 0.3}
			postEmbedding := []float32{0.15, 0.25, 0.35}
			expectedError := errors.New("chat service unavailable")

			redditResponse := &reddit.RedditResponse{
				Data: reddit.RedditData{
					Children: []reddit.RedditChild{
						{
							Data: reddit.RedditPostData{
								Title:       "Test Post",
								Selftext:    "Test content",
								URL:         "https://reddit.com/r/test/post",
								Score:       10,
								NumComments: 5,
								CreatedUTC:  float64(time.Now().Unix()),
								Permalink:   "/r/test/post",
							},
						},
					},
				},
			}

			mockLLMClient.EXPECT().GetEmbedding(ctx, "test topic").Return(topicEmbedding, nil)
			mockRedditService.EXPECT().SearchPosts("test", "test topic", 5).Return(redditResponse, nil)
			mockLLMClient.EXPECT().GetEmbedding(ctx, "Test Post. Test content").Return(postEmbedding, nil)
			mockLLMClient.EXPECT().Chat(ctx, mock.Anything).Return("", expectedError)

			// Act
			result, err := service.GetRelevantPosts(ctx, request)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "error getting relevance summary")
			assert.Empty(t, result.Posts)
		})
	})
}
