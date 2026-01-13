package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"google.golang.org/api/idtoken"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "news-subscriber-core/api/proto/v1"
)

// GoogleLogin handles the initial Google OAuth login
func (s *Service) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.GoogleLoginResponse, error) {
	// Verify the Google ID token
	payload, err := s.verifyGoogleToken(ctx, req.IdToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid google token: %v", err)
	}

	// Check if user exists by Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, payload.GoogleID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user: %v", err)
	}

	// User exists - log them in
	if user != nil {
		token, err := s.generateSessionToken()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
		}

		return &pb.GoogleLoginResponse{
			Status:       pb.LoginStatus_LOGGED_IN,
			SessionToken: token,
			UserId:       user.ID,
			Email:        user.Email,
		}, nil
	}

	// New user - needs invite code
	return &pb.GoogleLoginResponse{
		Status: pb.LoginStatus_NEEDS_INVITE,
		Email:  payload.Email,
	}, nil
}

// CompleteSignup completes the signup process with an invite code
func (s *Service) CompleteSignup(ctx context.Context, req *pb.CompleteSignupRequest) (*pb.CompleteSignupResponse, error) {
	// Verify the Google ID token
	payload, err := s.verifyGoogleToken(ctx, req.IdToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid google token: %v", err)
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByGoogleID(ctx, payload.GoogleID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user: %v", err)
	}
	if existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	// Validate invite code
	valid, err := s.inviteRepo.IsValid(ctx, req.InviteCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate invite code: %v", err)
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid or already used invite code")
	}

	// Create the user
	user, err := s.userRepo.Create(ctx, payload.Email, payload.GoogleID, payload.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Mark invite code as used
	if err := s.inviteRepo.MarkAsUsed(ctx, req.InviteCode, user.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark invite code as used: %v", err)
	}

	// Generate session token
	token, err := s.generateSessionToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.CompleteSignupResponse{
		SessionToken: token,
		UserId:       user.ID,
	}, nil
}

// Helper types and methods

type GoogleTokenPayload struct {
	Email    string
	GoogleID string
	Name     string
}

func (s *Service) verifyGoogleToken(ctx context.Context, idToken string) (*GoogleTokenPayload, error) {
	// Validate the ID token with Google
	payload, err := idtoken.Validate(ctx, idToken, s.cfg.Google.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Extract user information from claims
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email not found in token")
	}

	googleID, ok := payload.Claims["sub"].(string)
	if !ok || googleID == "" {
		return nil, fmt.Errorf("subject (user ID) not found in token")
	}

	name, ok := payload.Claims["name"].(string)
	if !ok || name == "" {
		// Name is optional, use email as fallback
		name = email
	}

	return &GoogleTokenPayload{
		Email:    email,
		GoogleID: googleID,
		Name:     name,
	}, nil
}

func (s *Service) generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
