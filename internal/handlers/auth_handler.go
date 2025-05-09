package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/iwtcode/user-order-api/internal/utils"
)

// Хэндлер для авторизации пользователей (REST API)
type AuthHandler struct {
	authService services.AuthService
}

// Структура запроса на вход
// Используется для валидации данных при логине
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Конструктор хэндлера авторизации
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// Валидация и разбор запроса
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Warn("Validation failed during login: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Вызов бизнес-логики авторизации
	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.Warn("Invalid credentials for email: %s", req.Email)
			status = http.StatusUnauthorized
			c.JSON(status, gin.H{"error": "Invalid email or password"})
			return
		}
		utils.Error("Login failed for email %s: %v", req.Email, err)
		c.JSON(status, gin.H{"error": "Login failed"})
		return
	}

	utils.Info("User logged in: %s", req.Email)
	// Формирование и отправка ответа
	c.JSON(http.StatusOK, gin.H{"token": token})
}
