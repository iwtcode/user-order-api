package models

import (
	"time"

	"gorm.io/gorm"
)

// Структура заказа для хранения в базе данных
type Order struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Product   string    `gorm:"type:varchar(255);not null" json:"product"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderCreateRequest содержит данные для создания заказа
// swagger:model
// Структура для запроса на создание заказа
type OrderCreateRequest struct {
	Product  string  `json:"product" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,gte=1"`
	Price    float64 `json:"price" binding:"required,gt=0"`
}

// OrderResponse содержит данные заказа
// swagger:model
// Структура для ответа API с данными заказа
type OrderResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// Вспомогательная функция для формирования ответа API по заказу
func BuildOrderResponse(order *Order) OrderResponse {
	return OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Product:   order.Product,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
	}
}
