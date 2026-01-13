package service

import (
	"context"
	"fmt"
	"time"

	pb "news-subscriber-core/api/proto/v1"
)

func (s *Service) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	fmt.Printf("SendMessage request: %v\n", req)
	return &pb.SendMessageResponse{
		Message: &pb.Message{
			Id:        "dummy_msg_id",
			Role:      "assistant",
			Content:   "This is a dummy response.",
			Timestamp: time.Now().Format(time.RFC3339),
		},
		SuggestedPreferences: []string{"Tech", "AI"},
	}, nil
}
