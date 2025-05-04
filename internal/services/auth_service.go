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
	if err != nil || user == nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
