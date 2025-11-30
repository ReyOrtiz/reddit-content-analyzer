package services

import (
	"time"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/contracts"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/reddit"
)

func MapRedditResponseToSubredditPostDto(
	post reddit.RedditChild,
	subredditName string,
	relevanceScore float64,
	isRelevant bool,
	relevanceSummary string,
) contracts.SubRedditPostDto {
	return contracts.SubRedditPostDto{
		SubredditName:    subredditName,
		Title:            post.Data.Title,
		Content:          post.Data.Selftext,
		Url:              post.Data.URL,
		Score:            post.Data.Score,
		NumComments:      post.Data.NumComments,
		CreatedAt:        time.Unix(int64(post.Data.CreatedUTC), 0),
		IsRelevant:       isRelevant,
		RelevanceScore:   relevanceScore,
		RelevanceSummary: relevanceSummary,
	}
}
