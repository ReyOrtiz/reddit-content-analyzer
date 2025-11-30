package services

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/contracts"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/llm"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/reddit"
)

type RelevanceService interface {
	GetRelevantPosts(ctx context.Context, request contracts.RelevanceRequestDto) (contracts.RelevanceResponseDto, error)
}

type relevanceService struct {
	logger        *zap.Logger
	llmClient     *llm.Client
	redditService RedditService
}

func NewRelevanceService() RelevanceService {
	redditService := NewRedditService()
	llmClient := llm.GetClient()
	return &relevanceService{
		logger:        logger.GetLogger(),
		llmClient:     llmClient,
		redditService: redditService,
	}
}

func (s *relevanceService) GetRelevantPosts(ctx context.Context, request contracts.RelevanceRequestDto) (contracts.RelevanceResponseDto, error) {
	s.logger.Info("Getting relevant posts", zap.Any("request", request))

	topicEmbedding, err := s.llmClient.GetEmbedding(ctx, request.Topic)
	if err != nil {
		return contracts.RelevanceResponseDto{}, errors.Wrap(err, "error getting topic embedding")
	}

	subredditPostDtos := make([]contracts.SubRedditPostDto, 0)
	for _, subreddit := range request.Subreddits {
		var subredditPosts *reddit.RedditResponse
		switch request.SearchMethod {
		case contracts.SearchMethodSearch:
			subredditPosts, err = s.redditService.SearchPosts(subreddit, request.Topic, request.Limit)
			if err != nil {
				return contracts.RelevanceResponseDto{}, errors.Wrap(err, "error getting subreddit posts")
			}
		case contracts.SearchMethodLatest:
			subredditPosts, err = s.redditService.GetPosts(subreddit, request.Limit)
			if err != nil {
				return contracts.RelevanceResponseDto{}, errors.Wrap(err, "error getting subreddit posts")
			}
		}

		evalSubredditPostDtos, err := s.evaluateSubredditPosts(
			ctx,
			subreddit,
			subredditPosts,
			request.Topic,
			topicEmbedding,
			request.RelevanceThreshold,
		)
		if err != nil {
			return contracts.RelevanceResponseDto{}, errors.Wrap(err, "error evaluating subreddit posts")
		}

		subredditPostDtos = append(subredditPostDtos, evalSubredditPostDtos...)
	}

	return contracts.RelevanceResponseDto{
		Posts: subredditPostDtos,
	}, nil
}

func (s *relevanceService) evaluateSubredditPosts(
	ctx context.Context,
	subredditName string,
	subredditPosts *reddit.RedditResponse,
	topic string,
	topicEmbedding []float32,
	relevanceThreshold float64,
) ([]contracts.SubRedditPostDto, error) {
	subredditPostDtos := make([]contracts.SubRedditPostDto, 0)

	for _, post := range subredditPosts.Data.Children {
		relevanceScore, err := s.getRelevanceScore(ctx, post.Data.Title, post.Data.Selftext, topicEmbedding)
		if err != nil {
			return nil, errors.Wrap(err, "error getting relevance score")
		}
		isRelevant := relevanceScore >= relevanceThreshold
		relevanceSummary, err := s.getRelevanceSummary(ctx, post.Data.Title, post.Data.Selftext, topic, relevanceThreshold, relevanceScore, isRelevant)
		if err != nil {
			return nil, errors.Wrap(err, "error getting relevance summary")
		}
		postDto := MapRedditResponseToSubredditPostDto(post, subredditName, relevanceScore, isRelevant, relevanceSummary)
		subredditPostDtos = append(subredditPostDtos, postDto)
	}
	return subredditPostDtos, nil
}

func (s *relevanceService) getRelevanceScore(ctx context.Context, title, content string, topicEmbedding []float32) (float64, error) {
	s.logger.Info("Getting relevance score",
		zap.String("title", title),
		zap.String("content", content),
	)

	text := fmt.Sprintf("%s. %s", title, content)
	embedding, err := s.llmClient.GetEmbedding(ctx, text)
	if err != nil {
		return 0, errors.Wrap(err, "error getting embedding")
	}
	cosineSimilarity := CosineSimilarity(embedding, topicEmbedding)

	s.logger.Info(
		"Relevance score calculated",
		zap.String("title", title),
		zap.Float64("cosine_similarity", cosineSimilarity),
	)
	return cosineSimilarity, nil
}

func (s *relevanceService) getRelevanceSummary(
	ctx context.Context,
	title, content, topic string,
	relevanceThreshold float64,
	relevanceScore float64,
	isRelevant bool,
) (string, error) {
	s.logger.Info("Getting relevance summary",
		zap.String("title", title),
		zap.String("content", content),
		zap.String("topic", topic),
		zap.Float64("relevance_score", relevanceScore),
	)

	prompt := fmt.Sprintf(`Given the following title, content, and topic, generate an explanation of the relevance of the content to the topic. The explanation should be a single sentence.
	
	# Topic: "%s"
	# Relevance Threshold: %f
	# Is Relevant: %t
	# Relevance Score: %f

	Reddit Post:

	# Title: "%s"
	# Content: 
	%s
	`, topic, relevanceThreshold, isRelevant, relevanceScore, title, content,
	)

	response, err := s.llmClient.Chat(ctx, []llm.Message{
		{
			Role:    "user",
			Content: prompt,
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "error getting chat response")
	}

	return response, nil
}
