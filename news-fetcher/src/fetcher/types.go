package fetcher

import (
	"context"
	"time"
)

// NewsArticle represents a unified news article structure
type NewsArticle struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"image_url,omitempty"`
	Source      string    `json:"source"`
	Author      string    `json:"author,omitempty"`
	PublishedAt time.Time `json:"published_at"`
	Category    string    `json:"category,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

// FetcherConfig holds configuration for news fetchers
type FetcherConfig struct {
	URL         string        `json:"url"`
	APIKey      string        `json:"api_key,omitempty"`
	RateLimit   time.Duration `json:"rate_limit"`
	Timeout     time.Duration `json:"timeout"`
	MaxRetries  int           `json:"max_retries"`
	UserAgent   string        `json:"user_agent"`
}

// NewsFetcher defines the interface that all news fetchers must implement
type NewsFetcher interface {
	// FetchNews retrieves news articles from the configured source
	FetchNews(ctx context.Context) ([]NewsArticle, error)
	
	// GetConfig returns the current configuration
	GetConfig() FetcherConfig
	
	// GetName returns the name of the fetcher
	GetName() string
	
	// Validate checks if the fetcher is properly configured
	Validate() error
}

// BaseFetcher provides common functionality for all news fetchers
type BaseFetcher struct {
	config FetcherConfig
	name   string
}

// NewBaseFetcher creates a new base fetcher with the given configuration
func NewBaseFetcher(name string, config FetcherConfig) *BaseFetcher {
	// Set default values if not provided
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.UserAgent == "" {
		config.UserAgent = "NewsFetcher/1.0"
	}
	
	return &BaseFetcher{
		config: config,
		name:   name,
	}
}

// GetConfig returns the fetcher configuration
func (b *BaseFetcher) GetConfig() FetcherConfig {
	return b.config
}

// GetName returns the fetcher name
func (b *BaseFetcher) GetName() string {
	return b.name
}

// Validate performs basic validation of the fetcher configuration
func (b *BaseFetcher) Validate() error {
	if b.config.URL == "" {
		return ErrInvalidURL
	}
	if b.config.RateLimit <= 0 {
		return ErrInvalidRateLimit
	}
	return nil
}

