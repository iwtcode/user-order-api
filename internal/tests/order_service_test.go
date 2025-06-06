package test

import (
	"context"
	"testing"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrderRepo struct {
	mock.Mock
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
func (m *mockOrderRepo) ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	args := m.Called(ctx, userID)
	orders, _ := args.Get(0).([]models.Order)
	return orders, args.Error(1)
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func TestOrderService_CreateOrder(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	userRepo := new(mockUserRepo)
	svc := services.NewOrderService(orderRepo, userRepo)
	ctx := context.Background()

	user := &models.User{Email: "a@b.com"}
	orderReq := &models.OrderCreateRequest{Product: "Book", Quantity: 2, Price: 10.5}

	userRepo.On("GetUserByID", ctx, uint(1)).Return(user, nil)
	orderRepo.On("CreateOrder", ctx, mock.AnythingOfType("*models.Order")).Return(nil)

	resultChan := svc.CreateOrder(ctx, 1, orderReq)
	result := <-resultChan
	assert.NoError(t, result.Err)
	assert.Equal(t, uint(1), result.Order.UserID)
	assert.Equal(t, orderReq.Product, result.Order.Product)
	assert.Equal(t, orderReq.Quantity, result.Order.Quantity)
	assert.Equal(t, orderReq.Price, result.Order.Price)
}

func TestOrderService_CreateOrder_UserNotFound(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	userRepo := new(mockUserRepo)
	svc := services.NewOrderService(orderRepo, userRepo)
	ctx := context.Background()

	orderReq := &models.OrderCreateRequest{Product: "Book", Quantity: 2, Price: 10.5}
	userRepo.On("GetUserByID", ctx, uint(2)).Return(nil, nil)

	resultChan := svc.CreateOrder(ctx, 2, orderReq)
	result := <-resultChan
	assert.ErrorIs(t, result.Err, services.ErrOrderUserNotFound)
	assert.Nil(t, result.Order)
}

func TestOrderService_ListOrdersByUserID(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	userRepo := new(mockUserRepo)
	svc := services.NewOrderService(orderRepo, userRepo)
	ctx := context.Background()

	user := &models.User{Email: "a@b.com"}
	orders := []models.Order{{Product: "Book"}, {Product: "Pen"}}
	userRepo.On("GetUserByID", ctx, uint(1)).Return(user, nil)
	orderRepo.On("ListOrdersByUserID", ctx, uint(1)).Return(orders, nil)

	result, err := svc.ListOrdersByUserID(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Book", result[0].Product)
}

func TestOrderService_ListOrdersByUserID_UserNotFound(t *testing.T) {
	orderRepo := new(mockOrderRepo)
	userRepo := new(mockUserRepo)
	svc := services.NewOrderService(orderRepo, userRepo)
	ctx := context.Background()

	userRepo.On("GetUserByID", ctx, uint(2)).Return(nil, nil)

	result, err := svc.ListOrdersByUserID(ctx, 2)
	assert.ErrorIs(t, err, services.ErrOrderUserNotFound)
	assert.Nil(t, result)
}
