package main

import (
	"log"

	"github.com/iwtcode/user-order-api/internal/config"
	"github.com/iwtcode/user-order-api/internal/handlers"
	"github.com/iwtcode/user-order-api/internal/middleware"
	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database Connection (GORM)
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Initialize Dependencies (Repository -> Service -> Handler)
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userRepo)
	// Initialize other repos, services, handlers here later...

	// 4. Initialize Gin Router
	router := gin.Default()

	// Auth route (no JWT required)
	router.POST("/auth/login", authHandler.Login)

	router.POST("/users", userHandler.CreateUser)

	// User routes (JWT required for all except creation)
	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.JWTAuthMiddleware())
	{
		userRoutes.GET("", userHandler.ListUsers)
		// Add other user routes (GET, PUT, DELETE) here later...
	}

	// 5. Start Server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
