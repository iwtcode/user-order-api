package services

import (
	"context"
	"fmt"

	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/utils"
)

// Интерфейс сервиса авторизации
type AuthService interface {
	// Выполняет вход пользователя по email и паролю, возвращает JWT-токен
	Login(ctx context.Context, email, password string) (string, error)
}

// Реализация сервиса авторизации
// Использует репозиторий пользователей для проверки данных
type authService struct {
	userRepo repository.UserRepository
}

// Конструктор сервиса авторизации
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Выполняет вход пользователя по email и паролю, возвращает JWT-токен
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("database error during login for email %s: %w", email, err)
	}
	// Проверяем пароль
	if user == nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}
	// Генерируем JWT-токен
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT for user id=%d: %w", user.ID, err)
	}
	return token, nil
}
