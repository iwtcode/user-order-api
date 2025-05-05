package services

import (
	"context"
	"errors"

	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		utils.Error("Database error during login for email %s: %v", email, err)
		return "", errors.New("database error: " + err.Error())
	}
	if user == nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		utils.Warn("Invalid credentials for email: %s", email)
		return "", ErrInvalidCredentials
	}
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.Error("Failed to generate JWT for user id=%d: %v", user.ID, err)
		return "", errors.New("failed to generate JWT: " + err.Error())
	}
	return token, nil
}
