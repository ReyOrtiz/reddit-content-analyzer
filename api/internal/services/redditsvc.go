package services

import (
	"go.uber.org/zap"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/reddit"
)

type RedditService interface {
	GetPosts(subreddit string, limit int) (*reddit.RedditResponse, error)
	SearchPosts(subreddit string, query string, limit int) (*reddit.RedditResponse, error)
}

type redditService struct {
	client reddit.Client
	logger *zap.Logger
}

func NewRedditService() RedditService {
	logger := logger.GetLogger()
	client := reddit.NewClient()
	return &redditService{
		client: *client,
		logger: logger,
	}
}

func (s *redditService) GetPosts(subreddit string, limit int) (*reddit.RedditResponse, error) {
	s.logger.Info(
		"Getting Reddit posts",
		zap.String("subreddit", subreddit),
		zap.Int("limit", limit),
	)

	posts, err := s.client.GetPosts(subreddit, limit)
	if err != nil {
		s.logger.Error("Error getting Reddit posts", zap.Error(err))
		return nil, err
	}

	s.logger.Info(
		"Reddit posts found",
		zap.Any("posts", posts),
		zap.Int("count", len(posts.Data.Children)),
	)
	return posts, nil
}

func (s *redditService) SearchPosts(subreddit string, query string, limit int) (*reddit.RedditResponse, error) {
	s.logger.Info(
		"Searching Reddit posts",
		zap.String("subreddit", subreddit),
		zap.String("query", query),
		zap.Int("limit", limit),
	)

	posts, err := s.client.SearchPosts(subreddit, query, limit)
	if err != nil {
		s.logger.Error("Error searching Reddit posts", zap.Error(err))
		return nil, err
	}

	s.logger.Info(
		"Reddit search results found",
		zap.Any("posts", posts),
		zap.Int("count", len(posts.Data.Children)),
	)
	return posts, nil
}
