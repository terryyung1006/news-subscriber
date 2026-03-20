package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"news-subscriber-core/tools"

	"github.com/google/uuid"
)

const FirstOnboardingMessage = `*Yawns and stretches circuits*

Oh hello! I just woke up as your personal news companion. I'm here to learn what matters to YOU and deliver news that's actually relevant to your life.

To get started, tell me: What topics or areas of news genuinely interest you? Are there specific industries, regions, or themes you care about?`

const onboardingSessionTTL = time.Hour // Session expires after 1 hour of inactivity

// OnboardingSession stored in Redis
type OnboardingSessionData struct {
	SessionID string              `json:"session_id"`
	UserID    string              `json:"user_id"`
	Messages  []OnboardingMessage `json:"messages"`
	CreatedAt string              `json:"created_at"`
}

type OnboardingMessage struct {
	Role      string `json:"role"` // "user" or "assistant"
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// OnboardingTask is sent to inference worker
type OnboardingTask struct {
	JobId   string              `json:"job_id"`
	Type    string              `json:"type"`
	Payload OnboardingPayload   `json:"payload"`
}

type OnboardingPayload struct {
	UserID       string              `json:"user_id"`
	UserName     string              `json:"user_name"`
	Conversation []OnboardingMessage `json:"conversation"`
}

// OnboardingResult from inference worker
type OnboardingResult struct {
	Response   string           `json:"response"`
	IsComplete bool             `json:"is_complete"`
	Memories   []ExtractedMemory `json:"memories"`
	Status     string           `json:"status"`
	Error      string           `json:"error,omitempty"`
}

type ExtractedMemory struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// HTTP request/response types
type StartOnboardingRequest struct {
	UserID string `json:"user_id"`
}

type StartOnboardingResponse struct {
	SessionID    string `json:"session_id"`
	FirstMessage string `json:"first_message"`
}

type SendOnboardingMessageRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type SendOnboardingMessageResponse struct {
	Response   string            `json:"response"`
	IsComplete bool              `json:"is_complete"`
	Memories   []ExtractedMemory `json:"memories,omitempty"`
}

type GetOnboardingStatusRequest struct {
	UserID string `json:"user_id"`
}

type GetOnboardingStatusResponse struct {
	NeedsOnboarding  bool   `json:"needs_onboarding"`
	HasActiveSession bool   `json:"has_active_session"`
	SessionID        string `json:"session_id,omitempty"`
}

type GetUserMemoriesRequest struct {
	UserID string `json:"user_id"`
}

type GetUserMemoriesResponse struct {
	Memories []MemoryCard `json:"memories"`
}

type MemoryCard struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
}

func getSessionKey(userID string) string {
	return fmt.Sprintf("onboarding_session:%s", userID)
}

func (s *Service) StartOnboarding(ctx context.Context, req *StartOnboardingRequest) (*StartOnboardingResponse, error) {
	sessionID := uuid.New().String()
	now := time.Now()

	// Create initial session with the first assistant message
	session := OnboardingSessionData{
		SessionID: sessionID,
		UserID:    req.UserID,
		Messages: []OnboardingMessage{
			{
				Role:      "assistant",
				Content:   FirstOnboardingMessage,
				Timestamp: now.Format(time.RFC3339),
			},
		},
		CreatedAt: now.Format(time.RFC3339),
	}

	// Store in Redis
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %v", err)
	}

	err = tools.RedisClient.Set(ctx, getSessionKey(req.UserID), sessionJSON, onboardingSessionTTL).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %v", err)
	}

	return &StartOnboardingResponse{
		SessionID:    sessionID,
		FirstMessage: FirstOnboardingMessage,
	}, nil
}

func (s *Service) SendOnboardingMessage(ctx context.Context, req *SendOnboardingMessageRequest) (*SendOnboardingMessageResponse, error) {
	// 1. Get current session from Redis
	sessionJSON, err := tools.RedisClient.Get(ctx, getSessionKey(req.UserID)).Result()
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %v", err)
	}

	var session OnboardingSessionData
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("failed to parse session: %v", err)
	}

	// 2. Add user message to conversation
	now := time.Now()
	session.Messages = append(session.Messages, OnboardingMessage{
		Role:      "user",
		Content:   req.Message,
		Timestamp: now.Format(time.RFC3339),
	})

	// 3. Get user name for personalization
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	userName := "User"
	if err == nil && user != nil {
		userName = user.Name
	}

	// 4. Send to inference worker
	jobID := uuid.New().String()
	task := OnboardingTask{
		JobId: jobID,
		Type:  "process_onboarding",
		Payload: OnboardingPayload{
			UserID:       req.UserID,
			UserName:     userName,
			Conversation: session.Messages,
		},
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %v", err)
	}

	err = tools.RedisClient.RPush(ctx, "inference_queue", taskJSON).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %v", err)
	}

	// 5. Poll for result (max 90 seconds)
	timeout := time.After(90 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	resultKey := fmt.Sprintf("job_result:%s", jobID)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("inference timed out")
		case <-ticker.C:
			val, err := tools.RedisClient.Get(ctx, resultKey).Result()
			if err == nil {
				var result OnboardingResult
				if err := json.Unmarshal([]byte(val), &result); err != nil {
					return nil, fmt.Errorf("failed to parse result: %v", err)
				}

				if result.Status == "failed" {
					return nil, fmt.Errorf("inference failed: %s", result.Error)
				}

				// 6. Add assistant response to session
				session.Messages = append(session.Messages, OnboardingMessage{
					Role:      "assistant",
					Content:   result.Response,
					Timestamp: time.Now().Format(time.RFC3339),
				})

				// 7. If complete, store memories and mark user as onboarded
				if result.IsComplete && len(result.Memories) >= 3 {
					// Store memories in database
					for _, mem := range result.Memories {
						_, err := s.memoryRepo.Create(ctx, req.UserID, mem.Title, mem.Content, mem.Category)
						if err != nil {
							return nil, fmt.Errorf("failed to store memory: %v", err)
						}
					}

					// Mark user as onboarded
					if err := s.userRepo.SetOnboardingCompleted(ctx, req.UserID); err != nil {
						return nil, fmt.Errorf("failed to update user: %v", err)
					}

					// Clean up session from Redis
					tools.RedisClient.Del(ctx, getSessionKey(req.UserID))

					return &SendOnboardingMessageResponse{
						Response:   result.Response,
						IsComplete: true,
						Memories:   result.Memories,
					}, nil
				}

				// 8. Update session in Redis (extend TTL)
				updatedSessionJSON, _ := json.Marshal(session)
				tools.RedisClient.Set(ctx, getSessionKey(req.UserID), updatedSessionJSON, onboardingSessionTTL)

				return &SendOnboardingMessageResponse{
					Response:   result.Response,
					IsComplete: false,
				}, nil
			}
		}
	}
}

func (s *Service) GetOnboardingStatus(ctx context.Context, req *GetOnboardingStatusRequest) (*GetOnboardingStatusResponse, error) {
	// Check if user has completed onboarding
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.OnboardingCompleted {
		return &GetOnboardingStatusResponse{
			NeedsOnboarding:  false,
			HasActiveSession: false,
		}, nil
	}

	// Check if there's an active session
	sessionJSON, err := tools.RedisClient.Get(ctx, getSessionKey(req.UserID)).Result()
	if err == nil {
		var session OnboardingSessionData
		if json.Unmarshal([]byte(sessionJSON), &session) == nil {
			return &GetOnboardingStatusResponse{
				NeedsOnboarding:  true,
				HasActiveSession: true,
				SessionID:        session.SessionID,
			}, nil
		}
	}

	return &GetOnboardingStatusResponse{
		NeedsOnboarding:  true,
		HasActiveSession: false,
	}, nil
}

func (s *Service) GetUserMemories(ctx context.Context, req *GetUserMemoriesRequest) (*GetUserMemoriesResponse, error) {
	memories, err := s.memoryRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories: %v", err)
	}

	cards := make([]MemoryCard, len(memories))
	for i, m := range memories {
		cards[i] = MemoryCard{
			ID:        m.ID,
			Title:     m.Title,
			Content:   m.Content,
			Category:  m.Category,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		}
	}

	return &GetUserMemoriesResponse{
		Memories: cards,
	}, nil
}
