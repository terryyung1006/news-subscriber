package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

// Producer handles Kafka message production
type Producer struct {
	producer *kafka.Producer
	topic    string
}

// KafkaMessage represents the unified message structure sent to Kafka
type KafkaMessage struct {
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
	FetchedAt   time.Time `json:"fetched_at"`
	FetcherName string    `json:"fetcher_name"`
}

// ProducerConfig holds configuration for the Kafka producer
type ProducerConfig struct {
	BootstrapServers string `json:"bootstrap_servers"`
	Topic           string `json:"topic"`
	SecurityProtocol string `json:"security_protocol,omitempty"`
	SASLMechanism   string `json:"sasl_mechanism,omitempty"`
	SASLUsername    string `json:"sasl_username,omitempty"`
	SASLPassword    string `json:"sasl_password,omitempty"`
	CompressionType string `json:"compression_type,omitempty"`
	Acks            string `json:"acks,omitempty"`
	Retries         int    `json:"retries,omitempty"`
}

// NewProducer creates a new Kafka producer
func NewProducer(config ProducerConfig) (*Producer, error) {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"client.id":         "news-fetcher",
		"acks":              "all",
		"retries":           3,
		"compression.type":  "snappy",
	}

	// Add optional configurations
	if config.SecurityProtocol != "" {
		kafkaConfig.SetKey("security.protocol", config.SecurityProtocol)
	}
	if config.SASLMechanism != "" {
		kafkaConfig.SetKey("sasl.mechanism", config.SASLMechanism)
	}
	if config.SASLUsername != "" {
		kafkaConfig.SetKey("sasl.username", config.SASLUsername)
	}
	if config.SASLPassword != "" {
		kafkaConfig.SetKey("sasl.password", config.SASLPassword)
	}
	if config.CompressionType != "" {
		kafkaConfig.SetKey("compression.type", config.CompressionType)
	}
	if config.Acks != "" {
		kafkaConfig.SetKey("acks", config.Acks)
	}
	if config.Retries > 0 {
		kafkaConfig.SetKey("retries", config.Retries)
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    config.Topic,
	}, nil
}

// PublishMessage publishes a message to Kafka
func (p *Producer) PublishMessage(ctx context.Context, message KafkaMessage) error {
	// Convert message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Kafka message
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(message.ID), // Use article ID as key for partitioning
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "source", Value: []byte(message.Source)},
			{Key: "fetcher", Value: []byte(message.FetcherName)},
			{Key: "category", Value: []byte(message.Category)},
		},
	}

	// Send message asynchronously
	deliveryChan := make(chan kafka.Event, 1)
	err = p.producer.Produce(kafkaMsg, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery confirmation or timeout
	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}
		logrus.Debugf("Message delivered to topic %s [%d] at offset %v",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("message delivery timeout")
	}

	return nil
}

// PublishMessages publishes multiple messages to Kafka
func (p *Producer) PublishMessages(ctx context.Context, messages []KafkaMessage) error {
	for _, message := range messages {
		if err := p.PublishMessage(ctx, message); err != nil {
			return fmt.Errorf("failed to publish message %s: %w", message.ID, err)
		}
	}
	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() {
	if p.producer != nil {
		p.producer.Flush(15 * 1000) // Wait up to 15 seconds for messages to be delivered
		p.producer.Close()
	}
}

// HealthCheck checks if the producer is healthy
func (p *Producer) HealthCheck(ctx context.Context) error {
	// Try to get metadata for the topic
	metadata, err := p.producer.GetMetadata(&p.topic, false, 5000)
	if err != nil {
		return fmt.Errorf("failed to get topic metadata: %w", err)
	}

	if len(metadata.Topics) == 0 {
		return fmt.Errorf("topic %s not found", p.topic)
	}

	return nil
}

