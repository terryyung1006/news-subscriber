package fetcher

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
)

const (
	NewsAPIFreeTierLimit = 100
)

// NewsAPIFetcher implements the NewsFetcher interface for NewsAPI.org
type NewsAPIFetcher struct {
	*BaseFetcher
	httpClient    *HTTPClient
	requestCount  int
	lastResetDate string
	mu            sync.Mutex
}

// NewsAPIResponse represents the response structure from NewsAPI
type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

// NewNewsAPIFetcher creates a new NewsAPI fetcher
func NewNewsAPIFetcher(config FetcherConfig) *NewsAPIFetcher {
	baseFetcher := NewBaseFetcher("newsapi", config)
	httpClient := NewHTTPClient(config.Timeout, config.RateLimit)
	
	return &NewsAPIFetcher{
		BaseFetcher: baseFetcher,
		httpClient:  httpClient,
	}
}

// FetchNews retrieves news from NewsAPI
func (n *NewsAPIFetcher) FetchNews(ctx context.Context) ([]NewsArticle, error) {
	if err := n.Validate(); err != nil {
		return nil, err
	}

	// Check and update daily quota
	remaining := n.trackRequest()
	if remaining < 0 {
		return nil, fmt.Errorf("daily request limit (%d) exceeded", NewsAPIFreeTierLimit)
	}
	log.Printf("[NewsAPI] Request %d/%d for today (%d remaining)",
		n.requestCount, NewsAPIFreeTierLimit, remaining)

	// Build the API URL with parameters
	apiURL, err := n.buildAPIURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build API URL: %w", err)
	}

	// Create HTTP request
	req, err := CreateRequest("GET", apiURL, n.config.UserAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key to headers
	req.Header.Set("X-API-Key", n.config.APIKey)

	// Make the request
	resp, err := n.httpClient.MakeRequest(ctx, req, n.config.MaxRetries)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Read response body
	body, err := ReadResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse the response
	var newsResponse NewsAPIResponse
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if newsResponse.Status != "ok" {
		return nil, fmt.Errorf("newsAPI returned error status: %s", newsResponse.Status)
	}

	// Convert to unified format
	articles := make([]NewsArticle, 0, len(newsResponse.Articles))
	for _, article := range newsResponse.Articles {
		articles = append(articles, NewsArticle{
			ID:          generateArticleID(article.URL),
			Title:       article.Title,
			Content:     article.Content,
			Summary:     article.Description,
			URL:         article.URL,
			ImageURL:    article.URLToImage,
			Source:      article.Source.Name,
			Author:      article.Author,
			PublishedAt: article.PublishedAt,
			Category:    "general", // NewsAPI doesn't provide category in this endpoint
		})
	}

	return articles, nil
}

// Validate checks if the NewsAPI fetcher is properly configured
func (n *NewsAPIFetcher) Validate() error {
	if err := n.BaseFetcher.Validate(); err != nil {
		return err
	}
	
	if n.config.APIKey == "" {
		return ErrAPIKeyRequired
	}
	
	return nil
}

// buildAPIURL constructs the NewsAPI URL with query parameters
func (n *NewsAPIFetcher) buildAPIURL() (string, error) {
	baseURL := n.config.URL
	if baseURL == "" {
		baseURL = "https://newsapi.org/v2/everything"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Add default query parameters
	params := u.Query()
	params.Set("apiKey", n.config.APIKey)
	params.Set("q", "business OR finance OR stocks")
	params.Set("language", "en")
	params.Set("sortBy", "publishedAt")
	params.Set("pageSize", "100")
	
	u.RawQuery = params.Encode()
	return u.String(), nil
}

// generateArticleID creates a unique ID for an article based on its URL
func generateArticleID(articleURL string) string {
	hash := sha256.Sum256([]byte(articleURL))
	return hex.EncodeToString(hash[:16])
}

// trackRequest tracks API requests and resets counter daily
// Returns remaining requests for the day (-1 if limit exceeded)
func (n *NewsAPIFetcher) trackRequest() int {
	n.mu.Lock()
	defer n.mu.Unlock()

	today := time.Now().Format("2006-01-02")

	// Reset counter if it's a new day
	if n.lastResetDate != today {
		n.requestCount = 0
		n.lastResetDate = today
		log.Printf("[NewsAPI] New day detected, resetting request counter")
	}

	n.requestCount++

	remaining := NewsAPIFreeTierLimit - n.requestCount
	if remaining < 20 {
		log.Printf("[NewsAPI] WARNING: Only %d requests remaining today", remaining)
	}

	return remaining
}
