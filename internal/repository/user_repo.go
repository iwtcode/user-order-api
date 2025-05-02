package repository

import (
	"errors"

	"github.com/iwtcode/user-order-api/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	ListUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error)
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
func (r *userRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user) // GORM automatically populates the ID field of the user struct
	return result.Error
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if user not found (not an application error)
		}
		return nil, result.Error // Return other DB errors
	}
	return &user, nil
}

func (r *userRepository) ListUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	query := r.db.Model(&models.User{})
	if minAge > 0 {
		query = query.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		query = query.Where("age <= ?", maxAge)
	}
	query.Count(&total)
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&users)
	return users, total, result.Error
}
