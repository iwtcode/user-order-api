package repository

import (
	"context"
	"errors"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/utils"
	"gorm.io/gorm"
)

// Интерфейс репозитория заказов для работы с БД
type OrderRepository interface {
	// Создаёт новый заказ в базе данных
	CreateOrder(ctx context.Context, order *models.Order) error
	// Возвращает список заказов пользователя по его ID
	ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error)
}

// Реализация репозитория заказов на GORM
type orderRepository struct {
	db *gorm.DB
}

// Конструктор репозитория заказов
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Создаёт новый заказ в базе данных
func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		utils.Error("Failed to create order in DB: %v", result.Error)
		return errors.New("failed to create order: " + result.Error.Error())
	}
	return nil
}

// Возвращает список заказов пользователя по его ID
func (r *orderRepository) ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	var orders []models.Order
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&orders)
	if result.Error != nil {
		utils.Error("Failed to list orders for user_id=%d: %v", userID, result.Error)
		return nil, result.Error
	}
	return orders, nil
}
