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

type InferenceTask struct {
	JobId   string      `json:"job_id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type QuestionPayload struct {
	UserId   string `json:"user_id"`
	Question string `json:"question"`
}

type InferenceResult struct {
	UserId   string `json:"user_id"`
	Question string `json:"question"`
	Response string `json:"response"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

func (s *Service) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// 1. Generate a Job ID
	jobID := uuid.New().String()

	// 2. Create the task payload
	task := InferenceTask{
		JobId: jobID,
		Type:  "process_question",
		Payload: QuestionPayload{
			UserId:   "user_123", // TODO: Get actual user ID from context
			Question: req.Message,
		},
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %v", err)
	}

	// 3. Push to Redis Queue
	// Using RPUSH to add to the tail of the list
	err = tools.RedisClient.RPush(ctx, "inference_queue", taskJSON).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %v", err)
	}

	// 4. Poll for result (Simple polling for now, max 30 seconds)
	// In a real production app, we might return the JobID to the client and let them poll,
	// or use a stream/websocket. For this "simple" version, we block and wait.
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	resultKey := fmt.Sprintf("job_result:%s", jobID)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("inference timed out")
		case <-ticker.C:
			// Check if result exists
			val, err := tools.RedisClient.Get(ctx, resultKey).Result()
			if err == nil {
				// Result found!
				var result InferenceResult
				if err := json.Unmarshal([]byte(val), &result); err != nil {
					return nil, fmt.Errorf("failed to parse result: %v", err)
				}

				if result.Status == "failed" {
					return nil, fmt.Errorf("inference failed: %s", result.Error)
				}

				return &pb.SendMessageResponse{
					Message: &pb.Message{
						Id:        jobID,
						Role:      "assistant",
						Content:   result.Response,
						Timestamp: time.Now().Format(time.RFC3339),
					},
				}, nil
			}
		}
	}
}
