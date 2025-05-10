package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/iwtcode/user-order-api/internal/utils"
)

// Хэндлер для работы с заказами (REST API)
type OrderHandler struct {
	orderService services.OrderService
}

// Конструктор хэндлера заказов
func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder godoc
// @Summary Создать заказ для пользователя
// @Description Создаёт новый заказ для пользователя по его ID. Пользователь может создавать заказы только для своего user_id.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body models.OrderCreateRequest true "Данные заказа"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]interface{}
// @Router /users/{id}/orders [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Получаем user_id из JWT (middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.Warn("User ID not found in token during order creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		utils.Error("Invalid user ID in token during order creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	// Проверяем, совпадает ли user_id из токена с id из path
	idParam := c.Param("id")
	var pathID uint
	_, err := fmt.Sscanf(idParam, "%d", &pathID)
	if err != nil || pathID == 0 {
		utils.Warn("Invalid user ID in path during order creation: %s", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in path"})
		return
	}
	if userID != pathID {
		utils.Warn("Access denied: user %d tried to create order for user %d", userID, pathID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: you can only operate with your own orders"})
		return
	}
	// Валидация и разбор запроса
	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Warn("Validation failed during order creation: %v", err)
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
	// Вызов бизнес-логики создания заказа (асинхронно)
	resultChan := h.orderService.CreateOrder(c.Request.Context(), userID, &req)
	result := <-resultChan
	if result.Err != nil {
		if errors.Is(result.Err, services.ErrOrderUserNotFound) {
			utils.Warn("Order creation failed: user not found (user_id=%d)", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		utils.Error("Failed to create order for user_id=%d: %v", userID, result.Err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	order := result.Order
	utils.Info("Order created: id=%d, user_id=%d, product=%s", order.ID, userID, order.Product)
	// Формирование и отправка ответа
	resp := models.BuildOrderResponse(order)
	c.JSON(http.StatusCreated, resp)
}

// GetOrdersByUserID godoc
// @Summary Получить заказы пользователя
// @Description Возвращает список заказов пользователя по его ID. Пользователь может просматривать только свои заказы.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/orders [get]
// @Security BearerAuth
func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	// Получаем user_id из JWT (middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.Warn("User ID not found in token during order list fetch")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		utils.Error("Invalid user ID in token during order list fetch")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	// Проверяем, совпадает ли user_id из токена с id из path
	idParam := c.Param("id")
	var pathID uint
	_, err := fmt.Sscanf(idParam, "%d", &pathID)
	if err != nil || pathID == 0 {
		utils.Warn("Invalid user ID in path during order list fetch: %s", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in path"})
		return
	}
	if userID != pathID {
		utils.Warn("Access denied: user %d tried to view orders for user %d", userID, pathID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: you can only view your own orders"})
		return
	}
	// Вызов бизнес-логики
	orders, err := h.orderService.ListOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrOrderUserNotFound) {
			utils.Warn("Order list fetch failed: user not found (user_id=%d)", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		utils.Error("Failed to fetch orders for user_id=%d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	utils.Info("Orders fetched for user_id=%d, count=%d", userID, len(orders))
	// Формирование и отправка ответа
	resp := make([]models.OrderResponse, len(orders))
	for i, o := range orders {
		resp[i] = models.BuildOrderResponse(&o)
	}
	c.JSON(http.StatusOK, resp)
}
