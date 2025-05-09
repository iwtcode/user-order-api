package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Загружаем .env один раз при инициализации пакета
func init() {
	err := godotenv.Load()
	if err != nil {
		Warn("No .env file found or error loading it, relying on environment variables. %v", err)
	}
}

// jwtSecret — секретный ключ для подписи токенов
var jwtSecret = []byte(getJWTSecret())

// Получает секретный ключ из .env или переменных окружения
func getJWTSecret() string {
	if value, exists := os.LookupEnv("JWT_SECRET"); exists {
		return value
	}
	return "your-secret-key"
}

// Получает время жизни токена из .env или переменных окружения
func getJWTExpiration() time.Duration {
	if value, exists := os.LookupEnv("JWT_EXPIRATION"); exists {
		dur, err := time.ParseDuration(value)
		if err == nil {
			return dur
		}
		Warn("Invalid JWT_EXPIRATION format: %s, using default 24h", value)
	}
	return 24 * time.Hour
}

// Генерирует JWT-токен для пользователя по его ID
func GenerateJWT(userID uint) (string, error) {
	expiration := getJWTExpiration()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Разбирает и валидирует JWT-токен, возвращает claims
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

// Возвращает ошибку истечения срока действия токена
func JwtErrTokenExpired() error {
	return jwt.ErrTokenExpired
}
