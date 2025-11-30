package reddit

// Reddit API response structures
type RedditResponse struct {
	Data RedditData `json:"data"`
}

type RedditData struct {
	Children []RedditChild `json:"children"`
}

type RedditChild struct {
	Data RedditPostData `json:"data"`
}

type RedditPostData struct {
	Title       string  `json:"title"`
	Selftext    string  `json:"selftext"`
	URL         string  `json:"url"`
	Score       int     `json:"score"`
	NumComments int     `json:"num_comments"`
	CreatedUTC  float64 `json:"created_utc"`
	Permalink   string  `json:"permalink"`
	Stickied    bool    `json:"stickied"` // Indicates if post is pinned/community highlight
}
