package services

import (
	"context"
	"fmt"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/repository"
)

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
		return nil, fmt.Errorf("failed to check user for order: %w", err)
	}
	if user == nil {
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
		return nil, fmt.Errorf("failed to create order in database for user_id=%d: %w", userID, err)
	}
	return order, nil
}

// Возвращает список заказов пользователя по его ID
func (s *orderService) ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user for orders: %w", err)
	}
	if user == nil {
		return nil, ErrOrderUserNotFound
	}
	// Получаем список заказов
	orders, err := s.orderRepo.ListOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for user_id=%d: %w", userID, err)
	}
	return orders, nil
}
