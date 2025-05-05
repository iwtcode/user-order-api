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
		utils.Error("Database error during login for email %s: %v", email, err)
		return "", errors.New("database error: " + err.Error())
	}
	// Проверяем пароль
	if user == nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		utils.Warn("Invalid credentials for email: %s", email)
		return "", ErrInvalidCredentials
	}
	// Генерируем JWT-токен
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.Error("Failed to generate JWT for user id=%d: %v", user.ID, err)
		return "", errors.New("failed to generate JWT: " + err.Error())
	}
	return token, nil
}
