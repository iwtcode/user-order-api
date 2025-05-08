package models

// Структура пользователя для хранения в базе данных
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Email        string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Age          int    `gorm:"not null" json:"age"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
}

// UserResponse содержит данные пользователя
// swagger:model
// Структура для ответа API с данными пользователя
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// CreateUserRequest содержит данные для создания пользователя
// swagger:model
// Структура для запроса на создание пользователя
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gte=1"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUserRequest содержит данные для обновления пользователя
// swagger:model
// Структура для запроса на обновление пользователя
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,gte=1"`
}

// Вспомогательная функция для формирования ответа API по пользователю
func BuildUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}
