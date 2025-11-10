package services

import (
	"context"
	"errors"
	"fmt"
	"user-service/internal/domain"
	"user-service/internal/dto"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtService *JWTService
	userRepo   repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, registerDto dto.RegisterRequest) (*domain.User, error) {
	existingUser, err := s.userRepo.FindByUsername(ctx, registerDto.Username)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user with that username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Username:     registerDto.Username,
		PasswordHash: string(hashedPassword),
		Balance:      0,
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = userID

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, userDto dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, userDto.Username)
	if err != nil {
		return &dto.TokenResponse{}, err
	}

	if user == nil {
		return &dto.TokenResponse{}, errors.New("bad login data")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDto.Password)); err != nil {
		return &dto.TokenResponse{}, errors.New("password is not correct")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return &dto.TokenResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refrseshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return &dto.TokenResponse{}, errors.New("failed to generate refresh token")
	}

	if err := s.userRepo.SaveRefreshToken(ctx, user.ID, refrseshToken); err != nil {
		return &dto.TokenResponse{}, errors.New("failed to save refresh token")
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refrseshToken,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60,
	}, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token invalid or expired")
	}

	user, err := s.userRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("user is not found")
	}

	if user == nil {
		return &dto.TokenResponse{}, errors.New("user is not found")
	}

	if user.ID != claims.UserID {
		return nil, errors.New("token does not belong to the user")
	}

	newAccessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	if err := s.userRepo.SaveRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, errors.New("failed to refresh ref-token")
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int) error {
	if err := s.userRepo.RevokeRefreshToken(ctx, userID); err != nil {
		return errors.New("logout error")
	}

	return nil
}
