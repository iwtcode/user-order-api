package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler holds the user service
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser handles the POST /users request
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			details := make([]string, 0, len(ve))
			for _, fe := range ve {
				details = append(details, fe.Field()+": "+fe.Tag())
			}
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation failed", "details": details})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		}
		return
	}

	newUser, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		// Логируем внутреннюю ошибку, но клиенту не раскрываем детали
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	response := models.BuildUserResponse(newUser)
	c.JSON(http.StatusCreated, response)
}

// ListUsers handles GET /users with pagination and filtering
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minAge, _ := strconv.Atoi(c.DefaultQuery("min_age", "0"))
	maxAge, _ := strconv.Atoi(c.DefaultQuery("max_age", "0"))

	if err1 != nil || err2 != nil || page < 1 || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page and limit must be positive integers"})
		return
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, limit, minAge, maxAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	respUsers := make([]models.UserResponse, len(users))
	for i, u := range users {
		respUsers[i] = models.BuildUserResponse(&u)
	}
	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"users": respUsers,
	})
}
