package service

import (
	"context"
	"fmt"

	"news-subscriber-core/src/config"
	"news-subscriber-core/src/repository/chroma"
	"news-subscriber-core/src/repository/postgres"

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
}

func NewService(cfg *config.Config, chroma *chroma.Client, userRepo *postgres.UserRepository, inviteRepo *postgres.InviteCodeRepository) *Service {
	return &Service{
		cfg:        cfg,
		chroma:     chroma,
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
	}
}

func (s *Service) Run(ctx context.Context) error {
	fmt.Println("Service running...")
	// Check ChromaDB connection
	hb, err := s.chroma.Heartbeat(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to ChromaDB: %v\n", err)
	} else {
		fmt.Printf("ChromaDB Heartbeat: %v\n", hb)
	}
	return nil
}
