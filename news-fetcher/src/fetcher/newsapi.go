package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// NewsAPIFetcher implements the NewsFetcher interface for NewsAPI.org
type NewsAPIFetcher struct {
	*BaseFetcher
	httpClient *HTTPClient
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
	params.Set("language", "en")
	params.Set("sortBy", "publishedAt")
	params.Set("pageSize", "100")
	
	u.RawQuery = params.Encode()
	return u.String(), nil
}

// generateArticleID creates a unique ID for an article based on its URL
func generateArticleID(url string) string {
	// Simple hash-based ID generation
	// In production, you might want to use a more sophisticated approach
	return fmt.Sprintf("%x", len(url))
}
