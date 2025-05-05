package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"
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
// @Description Создаёт новый заказ для пользователя по его ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body models.OrderCreateRequest true "Данные заказа"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]interface{}
// @Router /users/{id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Получение и проверка userID
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	// Валидация и разбор запроса
	var req models.OrderCreateRequest
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
	// Вызов бизнес-логики создания заказа
	order, err := h.orderService.CreateOrder(c.Request.Context(), uint(userID), &req)
	if err != nil {
		if errors.Is(err, services.ErrOrderUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	// Формирование и отправка ответа
	resp := models.BuildOrderResponse(order)
	c.JSON(http.StatusCreated, resp)
}

// GetOrdersByUserID godoc
// @Summary Получить заказы пользователя
// @Description Возвращает список заказов пользователя по его ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/orders [get]
func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	// Получение и проверка userID
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	// Вызов бизнес-логики
	orders, err := h.orderService.ListOrdersByUserID(c.Request.Context(), uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrOrderUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	// Формирование и отправка ответа
	resp := make([]models.OrderResponse, len(orders))
	for i, o := range orders {
		resp[i] = models.BuildOrderResponse(&o)
	}
	c.JSON(http.StatusOK, resp)
}
