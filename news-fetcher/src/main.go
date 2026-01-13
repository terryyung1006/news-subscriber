package main

import (
	"context"
	"news-fetcher/src/config"
	"news-fetcher/src/fetcher"
	"news-fetcher/src/kafka"
	"news-fetcher/src/worker"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	setupLogging(cfg.Logging)

	logrus.Info("Starting News Fetcher Worker")

	// Create Kafka producer
	kafkaProducer, err := kafka.NewProducer(kafka.ProducerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		Topic:            cfg.Kafka.Topic,
		SecurityProtocol: cfg.Kafka.SecurityProtocol,
		SASLMechanism:    cfg.Kafka.SASLMechanism,
		SASLUsername:     cfg.Kafka.SASLUsername,
		SASLPassword:     cfg.Kafka.SASLPassword,
		CompressionType:  cfg.Kafka.CompressionType,
		Acks:             cfg.Kafka.Acks,
		Retries:          cfg.Kafka.Retries,
	})
	if err != nil {
		logrus.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Test Kafka connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := kafkaProducer.HealthCheck(ctx); err != nil {
		logrus.Fatalf("Kafka health check failed: %v", err)
	}

	// Create fetcher registry
	registry := fetcher.NewFetcherRegistry()

	// Create worker
	newsWorker := worker.NewWorker(kafkaProducer, registry, cfg)

	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start worker in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		newsWorker.Start()
	}()

	// Wait for shutdown signal
	<-shutdownChan
	logrus.Info("Shutdown signal received, stopping worker...")

	// Stop worker gracefully
	newsWorker.Stop()
	wg.Wait()

	logrus.Info("News Fetcher Worker stopped")
}

// setupLogging configures the logging based on configuration
func setupLogging(cfg config.LoggingConfig) {
	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set log format
	if cfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}
