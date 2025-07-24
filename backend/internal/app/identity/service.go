package identity

import (
	"black-shop-service/internal/domain/identity"
	"black-shop-service/pkg/auth"
	"context"
	"go.uber.org/zap"
)

type Service struct {
	userRepo     identity.UserRepository
	logger       *zap.Logger
	tokenManager *auth.TokenManager
}

func NewService(userRepo identity.UserRepository, logger *zap.Logger, tokenManager *auth.TokenManager) *Service {
	return &Service{
		userRepo:     userRepo,
		logger:       logger.Named("identity_service"),
		tokenManager: tokenManager,
	}
}

func (s *Service) RegisterUser(ctx context.Context, name, email, plainPassword string) (*identity.User, error) {
	s.logger.Info("attempting to register new user", zap.String("email", email))

	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		s.logger.Warn("registration failed: email already exists", zap.String("email", email))
		return nil, identity.ErrEmailAlreadyExists
	}

	newUser, err := identity.NewUser(name, email, plainPassword)
	if err != nil {
		s.logger.Debug("failed to create new user object", zap.Error(err))
		return nil, err
	}

	if err := s.userRepo.Save(ctx, newUser); err != nil {
		s.logger.Debug("failed to save new user to repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user registered successfully", zap.String("user_id", string(newUser.ID)))
	return newUser, nil
}

func (s *Service) LoginUser(ctx context.Context, email, plainPassword string) (*identity.User, string, error) {
	s.logger.Info("attempting to login user", zap.String("email", email))

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login failed: user not found", zap.String("email", email))
		return nil, "", identity.ErrInvalidPassword
	}

	if err := user.CheckPassword(plainPassword); err != nil {
		s.logger.Warn("login failed: invalid password", zap.String("email", email))
		return nil, "", identity.ErrInvalidPassword
	}

	claims := auth.UserClaims{
		UserID: string(user.ID),
	}
	token, err := s.tokenManager.Generate(claims)
	if err != nil {
		s.logger.Error("failed to generate token for user", zap.String("user_id", string(user.ID)), zap.Error(err))
		return nil, "", err
	}

	s.logger.Info("user logged in successfully and token generated", zap.String("user_id", string(user.ID)))
	return user, token, nil
}
