package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}
func (m *mockUserService) ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	args := m.Called(ctx, page, limit, minAge, maxAge)
	users, _ := args.Get(0).([]models.User)
	total, _ := args.Get(1).(int64)
	return users, total, args.Error(2)
}
func (m *mockUserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	args := m.Called(ctx, userID)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}
func (m *mockUserService) UpdateUser(ctx context.Context, userID uint, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, userID, req)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}
func (m *mockUserService) DeleteUser(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		requestBody  gin.H
		mockSetup    func(m *mockUserService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:        "success",
			requestBody: gin.H{"email": "a@b.com", "password": "12345678", "name": "Test", "age": 20},
			mockSetup: func(m *mockUserService) {
				m.On("CreateUser", mock.Anything, &models.CreateUserRequest{Email: "a@b.com", Password: "12345678", Name: "Test", Age: 20}).Return(&models.User{Email: "a@b.com", Name: "Test", Age: 20}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{"email": "a@b.com", "name": "Test", "age": float64(20)},
		},
		{
			name:         "validation error",
			requestBody:  gin.H{"email": "bad", "password": "", "name": "", "age": 0},
			mockSetup:    func(m *mockUserService) {},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{"error": "Validation failed"},
		},
		{
			name:        "email exists",
			requestBody: gin.H{"email": "a@b.com", "password": "12345678", "name": "Test", "age": 20},
			mockSetup: func(m *mockUserService) {
				m.On("CreateUser", mock.Anything, &models.CreateUserRequest{Email: "a@b.com", Password: "12345678", Name: "Test", Age: 20}).Return(nil, services.ErrEmailExists)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Email already exists"},
		},
		{
			name:        "internal error",
			requestBody: gin.H{"email": "a@b.com", "password": "12345678", "name": "Test", "age": 20},
			mockSetup: func(m *mockUserService) {
				m.On("CreateUser", mock.Anything, &models.CreateUserRequest{Email: "a@b.com", Password: "12345678", Name: "Test", Age: 20}).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to create user"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := NewUserHandler(mockSvc)
			r := gin.Default()
			r.POST("/users", h.CreateUser)
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
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

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		query        string
		mockSetup    func(m *mockUserService)
		expectedCode int
		expectedLen  int
		expectedBody map[string]interface{}
	}{
		{
			name:  "success",
			query: "?page=1&limit=2",
			mockSetup: func(m *mockUserService) {
				users := []models.User{{Email: "a@b.com", Name: "A", Age: 20}, {Email: "b@b.com", Name: "B", Age: 21}}
				m.On("ListUsers", mock.Anything, 1, 2, 0, 0).Return(users, int64(2), nil)
			},
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name:         "bad params",
			query:        "?page=0&limit=0",
			mockSetup:    func(m *mockUserService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "page and limit must be positive integers"},
		},
		{
			name:  "internal error",
			query: "?page=1&limit=2",
			mockSetup: func(m *mockUserService) {
				m.On("ListUsers", mock.Anything, 1, 2, 0, 0).Return(nil, int64(0), errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to fetch users"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := NewUserHandler(mockSvc)
			r := gin.Default()
			r.GET("/users", h.ListUsers)
			req, _ := http.NewRequest(http.MethodGet, "/users"+tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusOK {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Equal(t, float64(tt.expectedLen), float64(len(resp["users"].([]interface{}))))
			} else if tt.expectedBody != nil {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				for k, v := range tt.expectedBody {
					assert.Equal(t, v, resp[k])
				}
			}
		})
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		userID       string
		mockSetup    func(m *mockUserService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:   "success",
			userID: "1",
			mockSetup: func(m *mockUserService) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(&models.User{Email: "a@b.com", Name: "A", Age: 20}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{"email": "a@b.com", "name": "A", "age": float64(20)},
		},
		{
			name:         "bad id",
			userID:       "abc",
			mockSetup:    func(m *mockUserService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Invalid user ID"},
		},
		{
			name:   "not found",
			userID: "2",
			mockSetup: func(m *mockUserService) {
				m.On("GetUserByID", mock.Anything, uint(2)).Return(nil, services.ErrUserNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{"error": "User not found"},
		},
		{
			name:   "internal error",
			userID: "1",
			mockSetup: func(m *mockUserService) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to fetch user"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := NewUserHandler(mockSvc)
			r := gin.Default()
			r.GET("/users/:id", h.GetUserByID)
			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
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

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		userID       string
		requestBody  gin.H
		mockSetup    func(m *mockUserService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:        "success",
			userID:      "1",
			requestBody: gin.H{"name": "NewName", "email": "new@b.com", "age": 22},
			mockSetup: func(m *mockUserService) {
				m.On("UpdateUser", mock.Anything, uint(1), &models.UpdateUserRequest{Name: "NewName", Email: "new@b.com", Age: 22}).Return(&models.User{Name: "NewName", Email: "new@b.com", Age: 22}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{"name": "NewName", "email": "new@b.com", "age": float64(22)},
		},
		{
			name:         "bad id",
			userID:       "abc",
			requestBody:  gin.H{"name": "NewName", "email": "new@b.com", "age": 22},
			mockSetup:    func(m *mockUserService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Invalid user ID"},
		},
		{
			name:         "validation error",
			userID:       "1",
			requestBody:  gin.H{"name": "", "email": "", "age": 0},
			mockSetup:    func(m *mockUserService) {},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{"error": "Validation failed"},
		},
		{
			name:        "not found",
			userID:      "2",
			requestBody: gin.H{"name": "NewName", "email": "new@b.com", "age": 22},
			mockSetup: func(m *mockUserService) {
				m.On("UpdateUser", mock.Anything, uint(2), &models.UpdateUserRequest{Name: "NewName", Email: "new@b.com", Age: 22}).Return(nil, services.ErrUserNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{"error": "User not found"},
		},
		{
			name:        "email exists",
			userID:      "1",
			requestBody: gin.H{"name": "NewName", "email": "new@b.com", "age": 22},
			mockSetup: func(m *mockUserService) {
				m.On("UpdateUser", mock.Anything, uint(1), &models.UpdateUserRequest{Name: "NewName", Email: "new@b.com", Age: 22}).Return(nil, services.ErrEmailExists)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Email already exists"},
		},
		{
			name:        "internal error",
			userID:      "1",
			requestBody: gin.H{"name": "NewName", "email": "new@b.com", "age": 22},
			mockSetup: func(m *mockUserService) {
				m.On("UpdateUser", mock.Anything, uint(1), &models.UpdateUserRequest{Name: "NewName", Email: "new@b.com", Age: 22}).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to update user"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := NewUserHandler(mockSvc)
			r := gin.Default()
			r.PUT("/users/:id", h.UpdateUser)
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(body))
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

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		userID       string
		mockSetup    func(m *mockUserService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:   "success",
			userID: "1",
			mockSetup: func(m *mockUserService) {
				m.On("DeleteUser", mock.Anything, uint(1)).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "bad id",
			userID:       "abc",
			mockSetup:    func(m *mockUserService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Invalid user ID"},
		},
		{
			name:   "not found",
			userID: "2",
			mockSetup: func(m *mockUserService) {
				m.On("DeleteUser", mock.Anything, uint(2)).Return(services.ErrUserNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{"error": "User not found"},
		},
		{
			name:   "internal error",
			userID: "1",
			mockSetup: func(m *mockUserService) {
				m.On("DeleteUser", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Failed to delete user"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := NewUserHandler(mockSvc)
			r := gin.Default()
			r.DELETE("/users/:id", h.DeleteUser)
			req, _ := http.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				for k, v := range tt.expectedBody {
					assert.Equal(t, v, resp[k])
				}
			}
		})
	}
}
