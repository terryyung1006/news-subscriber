package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Kafka   KafkaConfig    `mapstructure:"kafka"`
	Fetchers []FetcherConfig `mapstructure:"fetchers"`
	Logging LoggingConfig  `mapstructure:"logging"`
	Worker  WorkerConfig   `mapstructure:"worker"`
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	BootstrapServers string `mapstructure:"bootstrap_servers"`
	Topic           string `mapstructure:"topic"`
	SecurityProtocol string `mapstructure:"security_protocol"`
	SASLMechanism   string `mapstructure:"sasl_mechanism"`
	SASLUsername    string `mapstructure:"sasl_username"`
	SASLPassword    string `mapstructure:"sasl_password"`
	CompressionType string `mapstructure:"compression_type"`
	Acks            string `mapstructure:"acks"`
	Retries         int    `mapstructure:"retries"`
}

// FetcherConfig holds configuration for individual news fetchers
type FetcherConfig struct {
	Name      string        `mapstructure:"name"`
	Type      string        `mapstructure:"type"`
	URL       string        `mapstructure:"url"`
	APIKey    string        `mapstructure:"api_key"`
	RateLimit time.Duration `mapstructure:"rate_limit"`
	Timeout   time.Duration `mapstructure:"timeout"`
	MaxRetries int          `mapstructure:"max_retries"`
	UserAgent string        `mapstructure:"user_agent"`
	Enabled   bool          `mapstructure:"enabled"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// WorkerConfig holds worker-specific configuration
type WorkerConfig struct {
	FetchInterval time.Duration `mapstructure:"fetch_interval"`
	MaxWorkers    int           `mapstructure:"max_workers"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/news-fetcher")

	// Set default values
	setDefaults()

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("NEWS_FETCHER")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Kafka defaults
	viper.SetDefault("kafka.bootstrap_servers", "localhost:9092")
	viper.SetDefault("kafka.topic", "news")
	viper.SetDefault("kafka.compression_type", "snappy")
	viper.SetDefault("kafka.acks", "all")
	viper.SetDefault("kafka.retries", 3)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Worker defaults
	viper.SetDefault("worker.fetch_interval", "1h")
	viper.SetDefault("worker.max_workers", 5)
	viper.SetDefault("worker.graceful_shutdown_timeout", "30s")

	// Fetcher defaults
	viper.SetDefault("fetchers", []map[string]interface{}{
		{
			"name":         "newsapi",
			"type":         "newsapi",
			"url":          "https://newsapi.org/v2/everything",
			"rate_limit":   "1s",
			"timeout":      "30s",
			"max_retries":  3,
			"user_agent":   "NewsFetcher/1.0",
			"enabled":      true,
		},
		{
			"name":         "guardian",
			"type":         "guardian",
			"url":          "https://content.guardianapis.com/search",
			"rate_limit":   "1s",
			"timeout":      "30s",
			"max_retries":  3,
			"user_agent":   "NewsFetcher/1.0",
			"enabled":      true,
		},
	})
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.Kafka.BootstrapServers == "" {
		return fmt.Errorf("kafka bootstrap_servers is required")
	}

	if config.Kafka.Topic == "" {
		return fmt.Errorf("kafka topic is required")
	}

	if len(config.Fetchers) == 0 {
		return fmt.Errorf("at least one fetcher must be configured")
	}

	for i, fetcher := range config.Fetchers {
		if fetcher.Name == "" {
			return fmt.Errorf("fetcher[%d] name is required", i)
		}
		if fetcher.Type == "" {
			return fmt.Errorf("fetcher[%d] type is required", i)
		}
		if fetcher.URL == "" {
			return fmt.Errorf("fetcher[%d] url is required", i)
		}
	}

	return nil
}

