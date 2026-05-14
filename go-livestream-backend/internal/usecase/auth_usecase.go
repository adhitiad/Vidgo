package usecase

import (
	"context"
	"errors"

	"go-livestream-backend/internal/config"
	"go-livestream-backend/internal/domain"
	"go-livestream-backend/pkg/hash"
	"go-livestream-backend/pkg/jwt"
)

type authUsecase struct {
	userRepo domain.UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(userRepo domain.UserRepository, cfg *config.Config) domain.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Validate Role
	switch req.Role {
	case domain.RoleAdmin, domain.RoleAgency, domain.RoleHost, domain.RoleUser:
		// Valid
	default:
		return nil, errors.New("invalid role")
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !hash.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}
