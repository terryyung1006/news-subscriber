package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// GuardianFetcher implements the NewsFetcher interface for The Guardian API
type GuardianFetcher struct {
	*BaseFetcher
	httpClient *HTTPClient
}

// GuardianResponse represents the response structure from The Guardian API
type GuardianResponse struct {
	Response struct {
		Status      string `json:"status"`
		Total       int    `json:"total"`
		StartIndex  int    `json:"startIndex"`
		PageSize    int    `json:"pageSize"`
		CurrentPage int    `json:"currentPage"`
		Pages       int    `json:"pages"`
		Results     []struct {
			ID              string    `json:"id"`
			Type            string    `json:"type"`
			SectionID       string    `json:"sectionId"`
			SectionName     string    `json:"sectionName"`
			WebPublicationDate time.Time `json:"webPublicationDate"`
			WebTitle        string    `json:"webTitle"`
			WebURL          string    `json:"webUrl"`
			APIURL          string    `json:"apiUrl"`
			Fields          struct {
				Headline    string `json:"headline"`
				Body        string `json:"body"`
				Thumbnail   string `json:"thumbnail"`
				Byline      string `json:"byline"`
				TrailText   string `json:"trailText"`
			} `json:"fields"`
			Tags []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				WebTitle string `json:"webTitle"`
				WebURL   string `json:"webUrl"`
				APIURL   string `json:"apiUrl"`
			} `json:"tags"`
		} `json:"results"`
	} `json:"response"`
}

// NewGuardianFetcher creates a new Guardian fetcher
func NewGuardianFetcher(config FetcherConfig) *GuardianFetcher {
	baseFetcher := NewBaseFetcher("guardian", config)
	httpClient := NewHTTPClient(config.Timeout, config.RateLimit)
	
	return &GuardianFetcher{
		BaseFetcher: baseFetcher,
		httpClient:  httpClient,
	}
}

// FetchNews retrieves news from The Guardian API
func (g *GuardianFetcher) FetchNews(ctx context.Context) ([]NewsArticle, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}

	// Build the API URL with parameters
	apiURL, err := g.buildAPIURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build API URL: %w", err)
	}

	// Create HTTP request
	req, err := CreateRequest("GET", apiURL, g.config.UserAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make the request
	resp, err := g.httpClient.MakeRequest(ctx, req, g.config.MaxRetries)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Read response body
	body, err := ReadResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse the response
	var guardianResponse GuardianResponse
	if err := json.Unmarshal(body, &guardianResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if guardianResponse.Response.Status != "ok" {
		return nil, fmt.Errorf("guardian API returned error status: %s", guardianResponse.Response.Status)
	}

	// Convert to unified format
	articles := make([]NewsArticle, 0, len(guardianResponse.Response.Results))
	for _, article := range guardianResponse.Response.Results {
		// Extract tags
		tags := make([]string, 0, len(article.Tags))
		for _, tag := range article.Tags {
			tags = append(tags, tag.WebTitle)
		}

		articles = append(articles, NewsArticle{
			ID:          article.ID,
			Title:       article.WebTitle,
			Content:     article.Fields.Body,
			Summary:     article.Fields.TrailText,
			URL:         article.WebURL,
			ImageURL:    article.Fields.Thumbnail,
			Source:      "The Guardian",
			Author:      article.Fields.Byline,
			PublishedAt: article.WebPublicationDate,
			Category:    article.SectionName,
			Tags:        tags,
		})
	}

	return articles, nil
}

// Validate checks if the Guardian fetcher is properly configured
func (g *GuardianFetcher) Validate() error {
	if err := g.BaseFetcher.Validate(); err != nil {
		return err
	}
	
	if g.config.APIKey == "" {
		return ErrAPIKeyRequired
	}
	
	return nil
}

// buildAPIURL constructs the Guardian API URL with query parameters
func (g *GuardianFetcher) buildAPIURL() (string, error) {
	baseURL := g.config.URL
	if baseURL == "" {
		baseURL = "https://content.guardianapis.com/search"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Add query parameters
	params := u.Query()
	params.Set("api-key", g.config.APIKey)
	params.Set("show-fields", "headline,body,thumbnail,byline,trailText")
	params.Set("show-tags", "all")
	params.Set("page-size", "50")
	params.Set("order-by", "newest")
	
	u.RawQuery = params.Encode()
	return u.String(), nil
}
