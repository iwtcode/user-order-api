package main

import (
	"github.com/iwtcode/user-order-api/internal/config"
	"github.com/iwtcode/user-order-api/internal/handlers"
	"github.com/iwtcode/user-order-api/internal/middleware"
	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/iwtcode/user-order-api/internal/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/iwtcode/user-order-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен вместе с префиксом Bearer
// @scheme bearer
// @bearerFormat JWT

// Настраиваем маршруты HTTP API
func setupRoutes(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, orderHandler *handlers.OrderHandler) *gin.Engine {
	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/auth/login", authHandler.Login)
	router.POST("/users", userHandler.CreateUser)

	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.JWTAuthMiddleware())
	{
		userRoutes.GET("", userHandler.ListUsers)
		userRoutes.GET(":id", userHandler.GetUserByID)
		userRoutes.PUT(":id", userHandler.UpdateUser)
		userRoutes.DELETE(":id", userHandler.DeleteUser)
		userRoutes.POST(":id/orders", orderHandler.CreateOrder)
		userRoutes.GET(":id/orders", orderHandler.GetOrdersByUserID)
	}

	return router
}

// Главная функция запускает сервер приложения
func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Error("Failed to load configuration: %v", err)
		return
	}

	// Устанавливаем режим Gin
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Инициализация логгера
	utils.InitLogger(cfg.LogFile)

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString), &gorm.Config{
		Logger: &utils.GormLogger{},
	})
	if err != nil {
		utils.Error("Failed to connect to database: %v", err)
		return
	}

	// Инициализируем репозитории, сервисы и хэндлеры
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	userService := services.NewUserService(userRepo)
	orderService := services.NewOrderService(orderRepo, userRepo)
	authService := services.NewAuthService(userRepo)

	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	authHandler := handlers.NewAuthHandler(authService)

	// Настраиваем маршруты
	router := setupRoutes(userHandler, authHandler, orderHandler)

	// Запускаем сервер
	if err := router.Run(cfg.ServerPort); err != nil {
		utils.Error("Failed to start server: %v", err)
	}
}
