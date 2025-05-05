package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/utils"
)

var ErrOrderUserNotFound = errors.New("user not found for order")

// Интерфейс сервиса заказов, описывает бизнес-логику работы с заказами
type OrderService interface {
	// Создаёт новый заказ для пользователя
	CreateOrder(ctx context.Context, userID uint, req *models.OrderCreateRequest) (*models.Order, error)
	// Возвращает список заказов пользователя по его ID
	ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error)
}

// Реализация сервиса заказов
// Использует репозитории заказов и пользователей
type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

// Конструктор сервиса заказов
func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository) OrderService {
	return &orderService{orderRepo: orderRepo, userRepo: userRepo}
}

// Создаёт новый заказ для пользователя
func (s *orderService) CreateOrder(ctx context.Context, userID uint, req *models.OrderCreateRequest) (*models.Order, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		utils.Error("Failed to check user for order: %v", err)
		return nil, fmt.Errorf("failed to check user: %w", err)
	}
	if user == nil {
		utils.Warn("Attempt to create order for non-existent user: %d", userID)
		return nil, ErrOrderUserNotFound
	}
	// Формируем структуру заказа
	order := &models.Order{
		UserID:   userID,
		Product:  req.Product,
		Quantity: req.Quantity,
		Price:    req.Price,
	}
	// Сохраняем заказ в базе
	err = s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		utils.Error("Failed to create order in database for user_id=%d: %v", userID, err)
		return nil, fmt.Errorf("failed to create order in database: %w", err)
	}
	utils.Info("Order created: id=%d, user_id=%d, product=%s", order.ID, userID, order.Product)
	return order, nil
}

// Возвращает список заказов пользователя по его ID
func (s *orderService) ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrOrderUserNotFound
	}
	// Получаем список заказов
	orders, err := s.orderRepo.ListOrdersByUserID(ctx, userID)
	if err != nil {
		utils.Error("Failed to get orders for user_id=%d: %v", userID, err)
		return nil, err
	}
	return orders, nil
}
