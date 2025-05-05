package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iwtcode/user-order-api/internal/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Warn("Missing or invalid Authorization header: %s", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			if errors.Is(err, utils.JwtErrTokenExpired()) {
				utils.Warn("Token expired: %s", tokenString)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
				return
			}
			utils.Warn("Invalid token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			utils.Error("Invalid token payload: user_id missing. Claims: %+v", claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload: user_id missing"})
			return
		}
		utils.Info("Authenticated user_id: %d", uint(userID))
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
