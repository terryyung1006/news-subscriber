package service

import (
	"context"
	"fmt"
	"time"

	pb "news-subscriber-core/api/proto/v1"
)

func (s *Service) GetLatestReport(ctx context.Context, req *pb.GetLatestReportRequest) (*pb.GetLatestReportResponse, error) {
	fmt.Printf("GetLatestReport request: %v\n", req)
	return &pb.GetLatestReportResponse{
		Report: &pb.Report{
			Id:        "dummy_report_id",
			Content:   "# Dummy Report\n\nThis is a dummy report content.",
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (s *Service) GenerateReport(ctx context.Context, req *pb.GenerateReportRequest) (*pb.GenerateReportResponse, error) {
	fmt.Printf("GenerateReport request: %v\n", req)
	return &pb.GenerateReportResponse{
		Report: &pb.Report{
			Id:        "new_dummy_report_id",
			Content:   "# New Generated Report\n\nThis is a newly generated report.",
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}, nil
}
