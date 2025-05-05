package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Функции для хеширования и проверки паролей
// Используется bcrypt для безопасности

// Хеширует пароль пользователя
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверяет соответствие пароля и хеша
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
