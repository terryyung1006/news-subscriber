package tools

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	InferenceRequestsQueue = "inference.requests"
	InferenceResultsQueue  = "inference.results"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

var RabbitMQ *RabbitMQClient

func InitRabbitMQ(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	for _, queue := range []string{InferenceRequestsQueue, InferenceResultsQueue} {
		_, err = ch.QueueDeclare(
			queue,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("failed to declare queue %s: %w", queue, err)
		}
	}

	RabbitMQ = &RabbitMQClient{conn: conn, channel: ch}
	fmt.Printf("Connected to RabbitMQ at %s\n", url)
	return nil
}

func (r *RabbitMQClient) Publish(queue, correlationID string, body []byte) error {
	return r.channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			CorrelationId: correlationID,
			Body:          body,
		},
	)
}

func (r *RabbitMQClient) Consume(queue, consumerTag string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,
		consumerTag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
}

func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
