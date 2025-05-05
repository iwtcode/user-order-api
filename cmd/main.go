package main

import (
	"github.com/iwtcode/user-order-api/internal/config"
	"github.com/iwtcode/user-order-api/internal/handlers"
	"github.com/iwtcode/user-order-api/internal/middleware"
	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/iwtcode/user-order-api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupRoutes(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	return router
}

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Error("Failed to load configuration: %v", err)
		return
	}

	// 2. Initialize Database Connection (GORM)
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		utils.Error("Failed to connect to database: %v", err)
		return
	}

	// 3. Initialize Dependencies (Repository -> Service -> Handler)
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	// Initialize other repos, services, handlers here later...

	// 4. Initialize Gin Router
	router := setupRoutes(userHandler, authHandler)

	// 5. Start Server
	if err := router.Run(cfg.ServerPort); err != nil {
		utils.Error("Failed to start server: %v", err)
	}
}
