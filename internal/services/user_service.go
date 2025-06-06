package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/utils"

	"gorm.io/gorm"
)

var ErrEmailExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrOrderUserNotFound = errors.New("order user not found")

// Интерфейс сервиса пользователей, описывает бизнес-логику работы с пользователями
type UserService interface {
	// Создаёт нового пользователя
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	// Возвращает список пользователей с пагинацией и фильтрацией
	ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error)
	// Получает пользователя по ID
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	// Обновляет данные пользователя
	UpdateUser(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.User, error)
	// Удаляет пользователя по ID
	DeleteUser(ctx context.Context, id uint) error
}

// Реализация сервиса пользователей
// Использует репозиторий для доступа к данным
type userService struct {
	userRepo repository.UserRepository
}

// Конструктор сервиса пользователей
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Создаёт нового пользователя
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking for existing email: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("attempt to create user with duplicate email: %s %w", req.Email, ErrEmailExists)
	}
	// Хешируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password for email %s: %v", req.Email, err)
	}
	// Создаём пользователя в базе
	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in database for email %s: %v", req.Email, err)
	}
	// Возвращаем созданного пользователя
	return newUser, nil
}

// Возвращает список пользователей с пагинацией и фильтрацией
func (s *userService) ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	users, total, err := s.userRepo.ListUsers(ctx, page, limit, minAge, maxAge)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	return users, total, nil
}

// Получает пользователя по ID
func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Обновляет данные пользователя
func (s *userService) UpdateUser(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.User, error) {
	// Получаем пользователя по ID
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	// Проверяем уникальность email, если он изменился
	if user.Email != req.Email {
		existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailExists
		}
	}
	// Обновляем поля пользователя
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	// Сохраняем изменения
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Удаляет пользователя по ID
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	// Удаляем пользователя
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}
