package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPClient provides HTTP functionality with rate limiting and retries
type HTTPClient struct {
	client    *http.Client
	rateLimit time.Duration
	lastCall  time.Time
}

// NewHTTPClient creates a new HTTP client with rate limiting
func NewHTTPClient(timeout time.Duration, rateLimit time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		rateLimit: rateLimit,
	}
}

// MakeRequest makes an HTTP request with rate limiting and retries
func (c *HTTPClient) MakeRequest(ctx context.Context, req *http.Request, maxRetries int) (*http.Response, error) {
	// Apply rate limiting
	if err := c.applyRateLimit(); err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Create a new request with context for each attempt
		reqWithCtx := req.WithContext(ctx)
		
		resp, err := c.client.Do(reqWithCtx)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				logrus.Warnf("Request attempt %d failed: %v, retrying...", attempt+1, err)
				time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, err)
		}

		// Check for rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			if attempt < maxRetries {
				logrus.Warnf("Rate limited, waiting before retry %d", attempt+1)
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, ErrRateLimitExceeded
		}

		// Check for other HTTP errors
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			if attempt < maxRetries {
				logrus.Warnf("HTTP error %d, retrying...", resp.StatusCode)
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("HTTP error %d", resp.StatusCode)
		}

		return resp, nil
	}

	return nil, lastErr
}

// applyRateLimit ensures requests don't exceed the configured rate limit
func (c *HTTPClient) applyRateLimit() error {
	if c.rateLimit <= 0 {
		return nil
	}

	now := time.Now()
	timeSinceLastCall := now.Sub(c.lastCall)
	
	if timeSinceLastCall < c.rateLimit {
		sleepDuration := c.rateLimit - timeSinceLastCall
		logrus.Debugf("Rate limiting: sleeping for %v", sleepDuration)
		time.Sleep(sleepDuration)
	}
	
	c.lastCall = time.Now()
	return nil
}

// ReadResponseBody reads and closes the response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// CreateRequest creates an HTTP request with proper headers
func CreateRequest(method, url, userAgent string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

