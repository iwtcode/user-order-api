package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/utils"
)

// LoggerMiddleware логирует HTTP-запросы через кастомный логгер
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if errors != "" {
			utils.ErrorSrc(utils.GinSource, "%s | %s | %d | %s | %s | %s", method, path, status, clientIP, duration, errors)
		} else {
			utils.InfoSrc(utils.GinSource, "%s | %s | %d | %s | %s", method, path, status, clientIP, duration)
		}
	}
}
