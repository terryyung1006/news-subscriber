package fetcher

import "errors"

// Common errors for news fetchers
var (
	ErrInvalidURL       = errors.New("invalid URL provided")
	ErrInvalidRateLimit = errors.New("invalid rate limit provided")
	ErrAPIKeyRequired   = errors.New("API key is required for this fetcher")
	ErrInvalidResponse  = errors.New("invalid response from news source")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrTimeout          = errors.New("request timeout")
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")
)

