package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/handlers"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  gin.H
		mockSetup    func(m *mockAuthService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:        "success",
			requestBody: gin.H{"email": "test@example.com", "password": "pass123"},
			mockSetup: func(m *mockAuthService) {
				m.On("Login", mock.Anything, "test@example.com", "pass123").Return("token123", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{"token": "token123"},
		},
		{
			name:        "invalid credentials",
			requestBody: gin.H{"email": "test@example.com", "password": "wrong"},
			mockSetup: func(m *mockAuthService) {
				m.On("Login", mock.Anything, "test@example.com", "wrong").Return("", services.ErrInvalidCredentials)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{"error": "Invalid email or password"},
		},
		{
			name:         "validation error",
			requestBody:  gin.H{"email": "bademail", "password": ""},
			mockSetup:    func(m *mockAuthService) {},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{"error": "Validation failed"},
		},
		{
			name:        "internal error",
			requestBody: gin.H{"email": "test@example.com", "password": "pass123"},
			mockSetup: func(m *mockAuthService) {
				m.On("Login", mock.Anything, "test@example.com", "pass123").Return("", errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Login failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockAuthService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}
			h := handlers.NewAuthHandler(mockSvc)

			r := gin.Default()
			r.POST("/login", h.Login)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
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
