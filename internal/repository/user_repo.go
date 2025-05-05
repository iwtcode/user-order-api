package repository

import (
	"context"
	"errors"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/utils"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser inserts a new user record into the database
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		utils.Error("Failed to create user in DB: %v", result.Error)
		return errors.New("failed to create user: " + result.Error.Error())
	}
	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		utils.Error("Failed to get user by email: %v", result.Error)
		return nil, errors.New("failed to get user by email: " + result.Error.Error())
	}
	return &user, nil
}

func (r *userRepository) ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	query := r.db.WithContext(ctx).Model(&models.User{})
	if minAge > 0 {
		query = query.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		query = query.Where("age <= ?", maxAge)
	}
	if err := query.Count(&total).Error; err != nil {
		utils.Error("Failed to count users: %v", err)
		return nil, 0, errors.New("failed to count users: " + err.Error())
	}
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		utils.Error("Failed to list users: %v", result.Error)
		return nil, 0, errors.New("failed to list users: " + result.Error.Error())
	}
	return users, total, nil
}
