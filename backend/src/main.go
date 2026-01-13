package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"news-subscriber-core/src/config"
	"news-subscriber-core/src/repository/chroma"
	"news-subscriber-core/src/repository/postgres"
	"news-subscriber-core/src/service"

	pb "news-subscriber-core/api/proto/v1"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize PostgreSQL
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Run migrations automatically on startup
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	inviteRepo := postgres.NewInviteCodeRepository(db)

	// Initialize ChromaDB client
	chromaClient, err := chroma.NewClient(cfg.Chroma.Address)
	if err != nil {
		log.Fatalf("Failed to create ChromaDB client: %v", err)
	}

	// Initialize Service with all dependencies
	svc := service.NewService(cfg, chromaClient, userRepo, inviteRepo)

	// Run service (just a check for now)
	ctx := context.Background()
	go func() {
		if err := svc.Run(ctx); err != nil {
			log.Printf("Service error: %v", err)
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Register services here...
	pb.RegisterAuthServiceServer(s, svc)
	pb.RegisterChatServiceServer(s, svc)
	pb.RegisterReportServiceServer(s, svc)
	pb.RegisterUserServiceServer(s, svc)

	fmt.Printf("Server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
