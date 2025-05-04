package handlers

import (
	"log"
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
		log.Printf("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrInvalidCredentials {
			status = http.StatusUnauthorized
			log.Printf("Failed login attempt for email: %s", req.Email)
		} else {
			log.Printf("Error during login for email %s: %v", req.Email, err)
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User logged in: %s", req.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
