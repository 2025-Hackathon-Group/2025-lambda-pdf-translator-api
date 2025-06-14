package router

import (
	handler "2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/handlers"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/middleware"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/repository"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/service"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter initializes and returns a new Gin router
func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// Initialize services
	userService := service.NewUserService(os.Getenv("JWT_SECRET"))

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, userRepo, db)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(db, userService)

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}

		userRoutes := api.Group("/user").Use(authMiddleware)
		{
			userRoutes.GET("/", userHandler.GetUser)
			userRoutes.PUT("/", userHandler.UpdateUser)
			userRoutes.GET("/organizations", userHandler.GetUserOrganisations)
		}
	}

	return r
}
