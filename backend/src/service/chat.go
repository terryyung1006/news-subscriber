package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "news-subscriber-core/api/proto/v1"
	"news-subscriber-core/tools"

	"github.com/google/uuid"
)

type InferenceRequest struct {
	CorrelationID string            `json:"correlation_id"`
	UserID        string            `json:"user_id"`
	UserName      string            `json:"user_name"`
	Question      string            `json:"question"`
	Context       map[string]string `json:"context,omitempty"`
}

type InferenceResult struct {
	CorrelationID string `json:"correlation_id"`
	Answer        string `json:"answer"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
}

func (s *Service) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	correlationID := uuid.New().String()

	resultCh := make(chan InferenceResult, 1)

	s.pendingMu.Lock()
	s.pending[correlationID] = resultCh
	s.pendingMu.Unlock()

	defer func() {
		s.pendingMu.Lock()
		delete(s.pending, correlationID)
		s.pendingMu.Unlock()
	}()

	inferenceReq := InferenceRequest{
		CorrelationID: correlationID,
		UserID:        req.UserId,
		UserName:      req.UserName,
		Question:      req.Message,
	}

	body, err := json.Marshal(inferenceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inference request: %w", err)
	}

	if err := tools.RabbitMQ.Publish(tools.InferenceRequestsQueue, correlationID, body); err != nil {
		return nil, fmt.Errorf("failed to publish inference request: %w", err)
	}

	select {
	case result := <-resultCh:
		if result.Status == "failed" {
			return nil, fmt.Errorf("inference failed: %s", result.Error)
		}
		return &pb.SendMessageResponse{
			Message: &pb.Message{
				Id:        correlationID,
				Role:      "assistant",
				Content:   result.Answer,
				Timestamp: time.Now().Format(time.RFC3339),
			},
		}, nil
	case <-time.After(90 * time.Second):
		return nil, fmt.Errorf("inference timed out")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
