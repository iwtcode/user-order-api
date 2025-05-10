package services

import (
	"context"
	"fmt"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/repository"
)

// Интерфейс сервиса заказов, описывает бизнес-логику работы с заказами
type OrderService interface {
	// Создаёт новый заказ для пользователя (асинхронно)
	// Возвращает канал, в который будет отправлен результат
	CreateOrder(ctx context.Context, userID uint, req *models.OrderCreateRequest) <-chan OrderResult
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

// Создаёт новый заказ для пользователя (асинхронно)
func (s *orderService) CreateOrder(ctx context.Context, userID uint, req *models.OrderCreateRequest) <-chan OrderResult {
	resultChan := make(chan OrderResult, 1)
	go func() {
		// Проверяем, существует ли пользователь
		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			resultChan <- OrderResult{Order: nil, Err: fmt.Errorf("failed to check user for order: %w", err)}
			close(resultChan)
			return
		}
		if user == nil {
			resultChan <- OrderResult{Order: nil, Err: ErrOrderUserNotFound}
			close(resultChan)
			return
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
			resultChan <- OrderResult{Order: nil, Err: fmt.Errorf("failed to create order in database for user_id=%d: %w", userID, err)}
			close(resultChan)
			return
		}
		resultChan <- OrderResult{Order: order, Err: nil}
		close(resultChan)
	}()
	return resultChan
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

// Асинхронный результат создания заказа
// Используется для возврата результата из горутины
// Можно расширить при необходимости
type OrderResult struct {
	Order *models.Order
	Err   error
}
