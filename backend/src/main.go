package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"news-subscriber-core/src/config"
	"news-subscriber-core/src/repository/chroma"
	"news-subscriber-core/src/repository/postgres"
	"news-subscriber-core/src/service"
	"news-subscriber-core/tools"

	pb "news-subscriber-core/api/proto/v1"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Initialize Redis
	tools.InitRedis()

	// Initialize RabbitMQ
	if err := tools.InitRabbitMQ(cfg.RabbitMQ.URL); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer tools.RabbitMQ.Close()

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
	memoryRepo := postgres.NewMemoryRepository(db)

	// Initialize ChromaDB client
	chromaClient, err := chroma.NewClient(cfg.Chroma.Address)
	if err != nil {
		log.Fatalf("Failed to create ChromaDB client: %v", err)
	}

	// Initialize Service with all dependencies
	svc := service.NewService(cfg, chromaClient, userRepo, inviteRepo, memoryRepo)

	ctx := context.Background()
	if err := svc.Start(ctx); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

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

	// Start HTTP server for simple JSON API (temporary for easy frontend integration)
	go func() {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
			// Enable CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != "POST" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req struct {
				UserId   string `json:"user_id"`
				UserName string `json:"user_name"`
				Message  string `json:"message"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Call the service method
			// Note: We are adapting the HTTP request to the gRPC service method signature
			resp, err := svc.SendMessage(context.Background(), &pb.SendMessageRequest{
				UserId:   req.UserId,
				UserName: req.UserName,
				Message:  req.Message,
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		// Onboarding status endpoint
		httpMux.HandleFunc("/api/onboarding/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			var req service.GetOnboardingStatusRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			resp, err := svc.GetOnboardingStatus(context.Background(), &req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		// Start onboarding endpoint
		httpMux.HandleFunc("/api/onboarding/start", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			var req service.StartOnboardingRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			resp, err := svc.StartOnboarding(context.Background(), &req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		// Send onboarding message endpoint
		httpMux.HandleFunc("/api/onboarding/message", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			var req service.SendOnboardingMessageRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			resp, err := svc.SendOnboardingMessage(context.Background(), &req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		// Password login endpoint
		httpMux.HandleFunc("/api/auth/password-login", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != "POST" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req service.PasswordLoginRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			resp, err := svc.PasswordLogin(context.Background(), &req)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		// Get user memories endpoint
		httpMux.HandleFunc("/api/memories", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			userID := r.URL.Query().Get("userId")
			if userID == "" {
				http.Error(w, "userId query parameter required", http.StatusBadRequest)
				return
			}

			resp, err := svc.GetUserMemories(context.Background(), &service.GetUserMemoriesRequest{UserID: userID})
			if err != nil {
				http.Error(w, fmt.Sprintf("Service error: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		log.Println("HTTP Server listening at :8081")
		if err := http.ListenAndServe(":8081", httpMux); err != nil {
			log.Fatalf("failed to serve http: %v", err)
		}
	}()

	fmt.Printf("Server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
