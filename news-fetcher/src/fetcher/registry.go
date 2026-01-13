package fetcher

import (
	"fmt"
)

// FetcherRegistry manages available news fetchers
type FetcherRegistry struct {
	fetchers map[string]func(FetcherConfig) NewsFetcher
}

// NewFetcherRegistry creates a new fetcher registry
func NewFetcherRegistry() *FetcherRegistry {
	registry := &FetcherRegistry{
		fetchers: make(map[string]func(FetcherConfig) NewsFetcher),
	}
	
	// Register default fetchers
	registry.Register("newsapi", func(config FetcherConfig) NewsFetcher {
		return NewNewsAPIFetcher(config)
	})
	registry.Register("guardian", func(config FetcherConfig) NewsFetcher {
		return NewGuardianFetcher(config)
	})
	
	return registry
}

// Register adds a new fetcher factory to the registry
func (r *FetcherRegistry) Register(name string, factory func(FetcherConfig) NewsFetcher) {
	r.fetchers[name] = factory
}

// CreateFetcher creates a fetcher instance by name
func (r *FetcherRegistry) CreateFetcher(name string, config FetcherConfig) (NewsFetcher, error) {
	factory, exists := r.fetchers[name]
	if !exists {
		return nil, fmt.Errorf("unknown fetcher: %s", name)
	}
	
	return factory(config), nil
}

// GetAvailableFetchers returns a list of available fetcher names
func (r *FetcherRegistry) GetAvailableFetchers() []string {
	names := make([]string, 0, len(r.fetchers))
	for name := range r.fetchers {
		names = append(names, name)
	}
	return names
}
