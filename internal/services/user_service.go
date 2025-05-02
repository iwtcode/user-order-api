package services

import (
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
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	ListUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error)
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
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// 1. Check if email already exists
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Handle potential DB errors during the check
		return nil, fmt.Errorf("error checking for existing email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailExists // Return specific error for duplicate email
	}

	// 2. Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create the user model
	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}

	// 4. Save the user to the database
	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		// Could be a race condition or other DB error
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	// newUser.ID is now populated by GORM after successful creation
	return newUser, nil
}

func (s *userService) ListUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	return s.userRepo.ListUsers(page, limit, minAge, maxAge)
}
