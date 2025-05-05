package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/services"
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Вызов бизнес-логики авторизации
	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrInvalidCredentials {
			status = http.StatusUnauthorized
			c.JSON(status, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(status, gin.H{"error": "Login failed"})
		return
	}

	// Формирование и отправка ответа
	c.JSON(http.StatusOK, gin.H{"token": token})
}
