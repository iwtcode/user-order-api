package models

import "gorm.io/gorm"

// User represents the user model corresponding to the 'users' table
type User struct {
	gorm.Model          // Includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Email        string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Age          int    `gorm:"not null" json:"age"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"` // Exclude password hash from JSON responses
}

// UserResponse is the format for sending user data back in API responses (excluding sensitive info)
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// CreateUserRequest defines the expected structure for user creation JSON payload
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gte=1"`      // Example validation: age >= 1
	Password string `json:"password" binding:"required,min=8"` // Example validation: password min 8 chars
}

// BuildUserResponse creates a UserResponse from a User model
func BuildUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}
