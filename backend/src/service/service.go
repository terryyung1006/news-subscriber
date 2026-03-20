package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"news-subscriber-core/src/config"
	"news-subscriber-core/src/repository/chroma"
	"news-subscriber-core/src/repository/postgres"
	"news-subscriber-core/tools"

	pb "news-subscriber-core/api/proto/v1"
)

type Service struct {
	pb.UnimplementedAuthServiceServer
	pb.UnimplementedChatServiceServer
	pb.UnimplementedReportServiceServer
	pb.UnimplementedUserServiceServer

	cfg        *config.Config
	chroma     *chroma.Client
	userRepo   *postgres.UserRepository
	inviteRepo *postgres.InviteCodeRepository

	pendingMu sync.Mutex
	pending   map[string]chan InferenceResult
	memoryRepo *postgres.MemoryRepository
}

func NewService(cfg *config.Config, chroma *chroma.Client, userRepo *postgres.UserRepository, inviteRepo *postgres.InviteCodeRepository, memoryRepo *postgres.MemoryRepository) *Service {
	return &Service{
		cfg:        cfg,
		chroma:     chroma,
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
		pending:    make(map[string]chan InferenceResult),
		memoryRepo: memoryRepo,
	}
}

func (s *Service) Start(ctx context.Context) error {
	fmt.Println("Service starting...")

	// Check ChromaDB connection
	hb, err := s.chroma.Heartbeat(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to ChromaDB: %v\n", err)
	} else {
		fmt.Printf("ChromaDB Heartbeat: %v\n", hb)
	}

	deliveries, err := tools.RabbitMQ.Consume(tools.InferenceResultsQueue, "backend-results")
	if err != nil {
		return fmt.Errorf("failed to start consuming results: %w", err)
	}

	go s.consumeResults(ctx, deliveries)
	return nil
}

func (s *Service) consumeResults(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-deliveries:
			if !ok {
				return
			}
			var result InferenceResult
			if err := json.Unmarshal(msg.Body, &result); err != nil {
				fmt.Printf("Failed to unmarshal inference result: %v\n", err)
				msg.Nack(false, false)
				continue
			}

			s.pendingMu.Lock()
			ch, found := s.pending[result.CorrelationID]
			s.pendingMu.Unlock()

			if found {
				select {
				case ch <- result:
				default:
				}
			}

			msg.Ack(false)
		}
	}
}
