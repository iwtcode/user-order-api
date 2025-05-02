package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Use go get github.com/joho/godotenv
)

// Config holds application configuration
type Config struct {
	DBConnectionString string
	ServerPort         string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (useful for local development)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, relying on environment variables.")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "your-postgres-password")
	dbName := getEnv("DB_NAME", "userorderapi")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	serverPort := getEnv("SERVER_PORT", "8080")

	return &Config{
		DBConnectionString: dsn,
		ServerPort:         ":" + serverPort, // Prepend colon for ListenAndServe
	}, nil
}

// Helper function to get environment variables with a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
