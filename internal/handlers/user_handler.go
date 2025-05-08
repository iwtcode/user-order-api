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

// Хэндлер для работы с пользователями (REST API)
type UserHandler struct {
	userService services.UserService
}

// Конструктор хэндлера пользователей
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary Создать пользователя
// @Description Регистрирует нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]interface{}
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Валидация и разбор запроса
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

	// Вызов бизнес-логики создания пользователя
	newUser, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Формирование и отправка ответа
	response := models.BuildUserResponse(newUser)
	c.JSON(http.StatusCreated, response)
}

// ListUsers godoc
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с пагинацией и фильтрацией
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Размер страницы"
// @Param min_age query int false "Минимальный возраст"
// @Param max_age query int false "Максимальный возраст"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /users [get]
// @Security BearerAuth
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Получение параметров пагинации и фильтрации
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minAge, _ := strconv.Atoi(c.DefaultQuery("min_age", "0"))
	maxAge, _ := strconv.Atoi(c.DefaultQuery("max_age", "0"))

	if err1 != nil || err2 != nil || page < 1 || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page and limit must be positive integers"})
		return
	}

	// Вызов бизнес-логики
	users, total, err := h.userService.ListUsers(c.Request.Context(), page, limit, minAge, maxAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Формирование и отправка ответа
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

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Получение и проверка ID
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Вызов бизнес-логики
	user, err := h.userService.GetUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Формирование и отправка ответа
	response := models.BuildUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Обновить пользователя
// @Description Обновляет данные пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body models.UpdateUserRequest true "Данные для обновления"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]interface{}
// @Router /users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Получение и проверка ID
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Валидация и разбор запроса
	var req models.UpdateUserRequest
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

	// Вызов бизнес-логики
	user, err := h.userService.UpdateUser(c.Request.Context(), uint(userID), &req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Формирование и отправка ответа
	response := models.BuildUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 204 {string} string ""
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Получение и проверка ID
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Вызов бизнес-логики
	err = h.userService.DeleteUser(c.Request.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Формирование и отправка ответа
	c.Status(http.StatusNoContent)
}
