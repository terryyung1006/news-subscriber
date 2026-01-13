package worker

import (
	"context"
	"fmt"
	"news-fetcher/src/config"
	"news-fetcher/src/fetcher"
	"news-fetcher/src/kafka"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Worker manages the news fetching process
type Worker struct {
	producer *kafka.Producer
	registry *fetcher.FetcherRegistry
	config   *config.Config
	fetchers []fetcher.NewsFetcher
	stopChan chan struct{}
	stopped  bool
	mu       sync.RWMutex
}

// NewWorker creates a new worker instance
func NewWorker(producer *kafka.Producer, registry *fetcher.FetcherRegistry, config *config.Config) *Worker {
	return &Worker{
		producer: producer,
		registry: registry,
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// Start begins the worker process
func (w *Worker) Start() {
	logrus.Info("Worker started")

	// Initialize fetchers
	if err := w.initializeFetchers(); err != nil {
		logrus.Fatalf("Failed to initialize fetchers: %v", err)
	}

	// Start the main fetch loop
	ticker := time.NewTicker(w.config.Worker.FetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.fetchAllNews()
		case <-w.stopChan:
			logrus.Info("Worker stop signal received")
			return
		}
	}
}

// Stop gracefully stops the worker
func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.stopped {
		return
	}

	w.stopped = true
	close(w.stopChan)

	// Wait for graceful shutdown
	timeout := w.config.Worker.GracefulShutdownTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for any ongoing operations to complete
	select {
	case <-ctx.Done():
		logrus.Warn("Graceful shutdown timeout exceeded")
	case <-time.After(100 * time.Millisecond):
		logrus.Info("Worker stopped gracefully")
	}
}

// initializeFetchers creates and validates all configured fetchers
func (w *Worker) initializeFetchers() error {
	w.fetchers = make([]fetcher.NewsFetcher, 0, len(w.config.Fetchers))

	for _, fetcherConfig := range w.config.Fetchers {
		if !fetcherConfig.Enabled {
			logrus.Infof("Fetcher %s is disabled, skipping", fetcherConfig.Name)
			continue
		}

		// Convert config to fetcher config
		fetcherCfg := fetcher.FetcherConfig{
			URL:        fetcherConfig.URL,
			APIKey:     fetcherConfig.APIKey,
			RateLimit:  fetcherConfig.RateLimit,
			Timeout:    fetcherConfig.Timeout,
			MaxRetries: fetcherConfig.MaxRetries,
			UserAgent:  fetcherConfig.UserAgent,
		}

		// Create fetcher
		f, err := w.registry.CreateFetcher(fetcherConfig.Type, fetcherCfg)
		if err != nil {
			return fmt.Errorf("failed to create fetcher %s: %w", fetcherConfig.Name, err)
		}

		// Validate fetcher
		if err := f.Validate(); err != nil {
			return fmt.Errorf("fetcher %s validation failed: %w", fetcherConfig.Name, err)
		}

		w.fetchers = append(w.fetchers, f)
		logrus.Infof("Initialized fetcher: %s (%s)", fetcherConfig.Name, fetcherConfig.Type)
	}

	if len(w.fetchers) == 0 {
		return fmt.Errorf("no enabled fetchers configured")
	}

	return nil
}

// fetchAllNews fetches news from all configured sources
func (w *Worker) fetchAllNews() {
	logrus.Info("Starting news fetch cycle")

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, w.config.Worker.MaxWorkers)

	for _, f := range w.fetchers {
		wg.Add(1)
		go func(fetcher fetcher.NewsFetcher) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			w.fetchFromSource(fetcher)
		}(f)
	}

	wg.Wait()
	logrus.Info("Completed news fetch cycle")
}

// fetchFromSource fetches news from a specific source
func (w *Worker) fetchFromSource(f fetcher.NewsFetcher) {
	startTime := time.Now()
	fetcherName := f.GetName()

	logrus.Infof("Fetching news from %s", fetcherName)

	ctx, cancel := context.WithTimeout(context.Background(), f.GetConfig().Timeout)
	defer cancel()

	// Fetch articles
	articles, err := f.FetchNews(ctx)
	if err != nil {
		logrus.Errorf("Failed to fetch news from %s: %v", fetcherName, err)
		return
	}

	if len(articles) == 0 {
		logrus.Infof("No articles found from %s", fetcherName)
		return
	}

	// Convert to Kafka messages
	messages := kafka.ConvertArticlesToKafkaMessages(articles, fetcherName)

	// Publish to Kafka
	if err := w.producer.PublishMessages(ctx, messages); err != nil {
		logrus.Errorf("Failed to publish messages from %s: %v", fetcherName, err)
		return
	}

	duration := time.Since(startTime)
	logrus.Infof("Successfully processed %d articles from %s in %v", len(articles), fetcherName, duration)
}
