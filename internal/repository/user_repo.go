package repository

import (
	"context"
	"errors"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/utils"

	"gorm.io/gorm"
)

// Интерфейс репозитория пользователей для работы с БД
// Описывает методы для CRUD-операций и поиска пользователей
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
}

// Реализация репозитория пользователей на GORM
type userRepository struct {
	db *gorm.DB
}

// Конструктор репозитория пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Создаёт нового пользователя в базе данных
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		utils.Error("Failed to create user in DB: %v", result.Error)
		return errors.New("failed to create user: " + result.Error.Error())
	}
	return nil
}

// Получает пользователя по email
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

// Получает пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		utils.Error("Failed to get user by id: %v", result.Error)
		return nil, errors.New("failed to get user by id: " + result.Error.Error())
	}
	return &user, nil
}

// Возвращает список пользователей с пагинацией и фильтрацией по возрасту
func (r *userRepository) ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	// minAge/maxAge — фильтрация по возрасту
	// page/limit — постраничный вывод
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

// Обновляет данные пользователя
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		utils.Error("Failed to update user in DB: %v", result.Error)
		return errors.New("failed to update user: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Удаляет пользователя по ID
func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		utils.Error("Failed to delete user in DB: %v", result.Error)
		return errors.New("failed to delete user: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
