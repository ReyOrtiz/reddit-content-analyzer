package contracts

import "time"

type RelevanceResponseDto struct {
	Posts []SubRedditPostDto `json:"posts"`
}

type SubRedditPostDto struct {
	SubredditName    string    `json:"subreddit_name"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Url              string    `json:"url"`
	Score            int       `json:"score"`
	NumComments      int       `json:"num_comments"`
	CreatedAt        time.Time `json:"created_at"`
	IsRelevant       bool      `json:"is_relevant"`
	RelevanceScore   float64   `json:"relevance_score"`
	RelevanceSummary string    `json:"relevance_summary"`
}
