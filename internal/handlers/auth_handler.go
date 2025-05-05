package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrInvalidCredentials {
			status = http.StatusUnauthorized
			c.JSON(status, gin.H{"error": "Invalid email or password"})
			return
		}
		// Не раскрываем детали внутренней ошибки клиенту
		c.JSON(status, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
