package contracts

import "time"

type SearchMethod string

const (
	SearchMethodSearch SearchMethod = "search"
	SearchMethodLatest SearchMethod = "latest"
)

type RelevanceRequestDto struct {
	Topic              string       `json:"topic" binding:"required"`
	Subreddits         []string     `json:"subreddits"`
	RelevanceThreshold float64      `json:"relevance_threshold"`
	Limit              int          `json:"limit"`
	CreatedAfter       time.Time    `json:"created_after"`
	MinNumComments     int          `json:"min_num_comments"`
	SearchMethod       SearchMethod `json:"search_method" binding:"required,oneof=search latest"`
}
