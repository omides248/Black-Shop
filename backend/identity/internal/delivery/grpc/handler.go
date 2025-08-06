package grpc

import (
	pb "black-shop/api/proto/v1"
	"black-shop/internal/identity/domain"
	"black-shop/pkg/contextkeys"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedIdentityServiceServer
	service *Service
	logger  *zap.Logger
}

func NewGRPCServer(service *Service, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		service: service,
		logger:  logger.Named("identity_grpc_handler"),
	}
}

func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.logger.Info("received Register request", zap.String("email", req.GetEmail()))

	user, err := s.service.RegisterUser(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		s.logger.Error("failed to register user", zap.String("email", req.GetEmail()), zap.Error(err))

		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to register user")
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:    string(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Info("received Login request", zap.String("email", req.GetEmail()))

	user, token, err := s.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		s.logger.Warn("failed to login user", zap.String("email", req.GetEmail()), zap.Error(err))

		if errors.Is(err, domain.ErrInvalidPassword) || errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
		}
		return nil, status.Errorf(codes.Internal, "failed to login user")
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Id:    string(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
		Token: token,
	}, nil
}

func (s *GRPCServer) GetMyProfile(ctx context.Context, req *pb.GetMyProfileRequest) (*pb.GetMyProfileResponse, error) {
	s.logger.Info("received GetMyProfile request")

	userIDString, err := contextkeys.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Warn("GetMyProfile failed: userID not found in context")
		return nil, status.Errorf(codes.Unauthenticated, "user is not authenticated")
	}

	userID := domain.UserID(userIDString)

	user, err := s.service.GetUserProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user profile")
	}

	return &pb.GetMyProfileResponse{
		User: &pb.User{
			Id:    string(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
