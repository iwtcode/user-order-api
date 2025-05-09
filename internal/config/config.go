package config

import (
	"fmt"
	"os"

	"github.com/iwtcode/user-order-api/internal/utils"
	"github.com/joho/godotenv"
)

// Структура для хранения конфигурации приложения (строка подключения к БД и порт сервера)
type Config struct {
	DBConnectionString string
	ServerPort         string
	GinMode            string
	LogFile            string
}

// Функция загружает конфигурацию из .env файла или переменных окружения
func LoadConfig() (*Config, error) {
	// Пытаемся загрузить .env файл
	err := godotenv.Load()
	if err != nil {
		utils.Warn("No .env file found or error loading it, relying on environment variables. %v", err)
	}

	// Формируем строку подключения к базе данных
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "db-password")
	dbName := getEnv("DB_NAME", "user_order_api")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	serverPort := getEnv("SERVER_PORT", "8080")
	ginMode := getEnv("GIN_MODE", "debug")
	logFile := getEnv("LOG_FILE", "")

	// Возвращаем структуру конфигурации
	return &Config{
		DBConnectionString: dsn,
		ServerPort:         ":" + serverPort,
		GinMode:            ginMode,
		LogFile:            logFile,
	}, nil
}

// Вспомогательная функция для получения переменной окружения с дефолтным значением
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
