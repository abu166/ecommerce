package user

import (
	"context"
	"ecommerce/internal/user/application"
	"ecommerce/internal/user/domain"
	"ecommerce/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	svc *application.Service
}

func NewServer(svc *application.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.UserResponse, error) {
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username, password, and email are required")
	}
	u := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := s.svc.Register(ctx, u); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	return &proto.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (s *Server) AuthenticateUser(ctx context.Context, req *proto.AuthenticateUserRequest) (*proto.AuthResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password are required")
	}
	token, err := s.svc.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	return &proto.AuthResponse{Token: token}, nil
}

func (s *Server) GetUserProfile(ctx context.Context, req *proto.GetUserProfileRequest) (*proto.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}
	u, err := s.svc.GetProfile(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return &proto.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}
