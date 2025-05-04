package handlers

import (
	"errors"
	"log"
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

	// Bind JSON request body to the CreateUserRequest struct
	// Use ShouldBindJSON for better error handling (doesn't abort)
	if err := c.ShouldBindJSON(&req); err != nil {
		// Provide more specific validation error messages
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			// Customize error response based on validation errors
			// For simplicity, just returning the first error for now
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": ve.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		}
		return
	}

	// Call the service to create the user
	newUser, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			log.Printf("Attempt to create user with existing email: %s", req.Email)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	log.Printf("User created: id=%d, email=%s", newUser.ID, newUser.Email)
	response := models.BuildUserResponse(newUser)
	c.JSON(http.StatusCreated, response)
}

// ListUsers handles GET /users with pagination and filtering
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minAge, _ := strconv.Atoi(c.DefaultQuery("min_age", "0"))
	maxAge, _ := strconv.Atoi(c.DefaultQuery("max_age", "0"))

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, limit, minAge, maxAge)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	log.Printf("Fetched users: page=%d, limit=%d, total=%d", page, limit, total)
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
