package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents a Reddit API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// NewClient creates a new Reddit client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   "https://www.reddit.com",
		userAgent: "reddit-content-analyzer/1.0",
	}
}

// GetPosts retrieves a list of posts from a given subreddit
// limit specifies the maximum number of posts to retrieve (default: 25, max: 100)
func (c *Client) GetPosts(subreddit string, limit int) (*RedditResponse, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	url := fmt.Sprintf("%s/r/%s/.json?limit=%d", c.baseURL, subreddit, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reddit API returned status %d: %s", resp.StatusCode, string(body))

	}

	var redditResponse *RedditResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return redditResponse, nil
}

// SearchPosts searches for posts in a subreddit by query terms
// limit specifies the maximum number of posts to retrieve (default: 25, max: 100)
func (c *Client) SearchPosts(subreddit string, query string, limit int) (*RedditResponse, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	// Reddit search endpoint with restrict_sr=true to limit search to the subreddit
	// URL encode the query parameter
	encodedQuery := url.QueryEscape(query)
	url := fmt.Sprintf("%s/r/%s/search.json?q=%s&restrict_sr=true&limit=%d", c.baseURL, subreddit, encodedQuery, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reddit API returned status %d: %s", resp.StatusCode, string(body))

	}

	var redditResponse *RedditResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return redditResponse, nil
}
