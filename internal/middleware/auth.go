package middleware

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := userService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", claims).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
