package service

import (
	"context"
	"fmt"

	pb "news-subscriber-core/api/proto/v1"
)

func (s *Service) GetPreferences(ctx context.Context, req *pb.GetPreferencesRequest) (*pb.GetPreferencesResponse, error) {
	fmt.Printf("GetPreferences request: %v\n", req)
	return &pb.GetPreferencesResponse{
		Preferences: []*pb.Preference{
			{Id: "1", Content: "Tech"},
			{Id: "2", Content: "Finance"},
		},
	}, nil
}

func (s *Service) AddPreference(ctx context.Context, req *pb.AddPreferenceRequest) (*pb.AddPreferenceResponse, error) {
	fmt.Printf("AddPreference request: %v\n", req)
	return &pb.AddPreferenceResponse{
		Preference: &pb.Preference{
			Id:      "new_id",
			Content: req.Content,
		},
	}, nil
}

func (s *Service) UpdatePreference(ctx context.Context, req *pb.UpdatePreferenceRequest) (*pb.UpdatePreferenceResponse, error) {
	fmt.Printf("UpdatePreference request: %v\n", req)
	return &pb.UpdatePreferenceResponse{
		Preference: &pb.Preference{
			Id:      req.PreferenceId,
			Content: req.Content,
		},
	}, nil
}

func (s *Service) DeletePreference(ctx context.Context, req *pb.DeletePreferenceRequest) (*pb.DeletePreferenceResponse, error) {
	fmt.Printf("DeletePreference request: %v\n", req)
	return &pb.DeletePreferenceResponse{
		Success: true,
	}, nil
}
