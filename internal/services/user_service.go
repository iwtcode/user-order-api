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

// Custom error for duplicate email
var ErrEmailExists = errors.New("email already exists")

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error)
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser handles the logic for creating a new user
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking for existing email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}
	return newUser, nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	return s.userRepo.ListUsers(ctx, page, limit, minAge, maxAge)
}
