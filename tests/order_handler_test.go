package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/handlers"
	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrderService struct {
	mock.Mock
}

func (m *mockOrderService) CreateOrder(ctx context.Context, userID uint, req *models.OrderCreateRequest) (*models.Order, error) {
	args := m.Called(ctx, userID, req)
	order, _ := args.Get(0).(*models.Order)
	return order, args.Error(1)
}

func (m *mockOrderService) ListOrdersByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	args := m.Called(ctx, userID)
	orders, _ := args.Get(0).([]models.Order)
	return orders, args.Error(1)
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		userID       string
		requestBody  gin.H
		mockSetup    func(m *mockOrderService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:        "success",
			userID:      "1",
			requestBody: gin.H{"product": "Book", "quantity": 2, "price": 10.5},
			mockSetup: func(m *mockOrderService) {
				m.On("CreateOrder", mock.Anything, uint(1), &models.OrderCreateRequest{Product: "Book", Quantity: 2, Price: 10.5}).Return(&models.Order{ID: 1, UserID: 1, Product: "Book", Quantity: 2, Price: 10.5, CreatedAt: time.Now()}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{"user_id": float64(1), "product": "Book", "quantity": float64(2), "price": 10.5},
		},
		{
			name:         "invalid user id",
			userID:       "abc",
			requestBody:  gin.H{"product": "Book", "quantity": 2, "price": 10.5},
			mockSetup:    func(m *mockOrderService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Invalid user_id"},
		},
		{
			name:         "validation error",
			userID:       "1",
			requestBody:  gin.H{"product": "", "quantity": 0, "price": 0},
			mockSetup:    func(m *mockOrderService) {},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{"error": "Validation failed"},
		},
		{
			name:        "user not found",
			userID:      "2",
			requestBody: gin.H{"product": "Book", "quantity": 2, "price": 10.5},
			mockSetup: func(m *mockOrderService) {
				m.On("CreateOrder", mock.Anything, uint(2), &models.OrderCreateRequest{Product: "Book", Quantity: 2, Price: 10.5}).Return(nil, services.ErrOrderUserNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{"error": "User not found"},
		},
		{
			name:        "internal error",
			userID:      "1",
			requestBody: gin.H{"product": "Book", "quantity": 2, "price": 10.5},
			mockSetup: func(m *mockOrderService) {
				m.On("CreateOrder", mock.Anything, uint(1), &models.OrderCreateRequest{Product: "Book", Quantity: 2, Price: 10.5}).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to create order"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockOrderService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := handlers.NewOrderHandler(mockSvc)

			r := gin.Default()
			r.POST("/users/:id/orders", h.CreateOrder)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/users/"+tt.userID+"/orders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			for k, v := range tt.expectedBody {
				assert.Equal(t, v, resp[k])
			}
		})
	}
}

func TestOrderHandler_GetOrdersByUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		userID       string
		mockSetup    func(m *mockOrderService)
		expectedCode int
		expectedLen  int
		expectedBody interface{}
	}{
		{
			name:   "success",
			userID: "1",
			mockSetup: func(m *mockOrderService) {
				orders := []models.Order{{ID: 1, UserID: 1, Product: "Book", Quantity: 2, Price: 10.5, CreatedAt: time.Now()}}
				m.On("ListOrdersByUserID", mock.Anything, uint(1)).Return(orders, nil)
			},
			expectedCode: http.StatusOK,
			expectedLen:  1,
		},
		{
			name:         "invalid user id",
			userID:       "abc",
			mockSetup:    func(m *mockOrderService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Invalid user_id"},
		},
		{
			name:   "user not found",
			userID: "2",
			mockSetup: func(m *mockOrderService) {
				m.On("ListOrdersByUserID", mock.Anything, uint(2)).Return(nil, services.ErrOrderUserNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{"error": "User not found"},
		},
		{
			name:   "internal error",
			userID: "1",
			mockSetup: func(m *mockOrderService) {
				m.On("ListOrdersByUserID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to fetch orders"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockOrderService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := handlers.NewOrderHandler(mockSvc)

			r := gin.Default()
			r.GET("/users/:id/orders", h.GetOrdersByUserID)

			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.userID+"/orders", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusOK {
				var resp []map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Equal(t, tt.expectedLen, len(resp))
			} else if tt.expectedBody != nil {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				for k, v := range tt.expectedBody.(map[string]interface{}) {
					assert.Equal(t, v, resp[k])
				}
			}
		})
	}
}
