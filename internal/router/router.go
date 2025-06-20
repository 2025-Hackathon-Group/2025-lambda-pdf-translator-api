package router

import (
	handler "2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/handlers"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/middleware"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/repository"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/service"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "2025-Hackathon-Group/2025-lambda-pdf-translator-api/docs"
)

// NewRouter initializes and returns a new Gin router
func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// Initialize services
	userService := service.NewUserService(os.Getenv("JWT_SECRET"))

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bucketRepository := repository.NewMinioBucketRepository(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"))

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, userRepo, db)
	bucketHandler := handler.NewBucketHandler(bucketRepository, repository.NewFileRepository(db))

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(db, userService)

	api := r.Group("")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

		userRoutes := api.Group("/me").Use(authMiddleware)
		{
			userRoutes.GET("", userHandler.GetUser)
			userRoutes.PATCH("", userHandler.UpdateUser)
			userRoutes.GET("/organisations", userHandler.GetUserOrganisations)
		}
		fileRoutes := api.Group("/files")
		{
			fileRoutes.Use(authMiddleware)
			fileRoutes.POST("/", bucketHandler.UploadFile)
			fileRoutes.GET("/:file_id", bucketHandler.GetFileByID)
			fileRoutes.GET("/:file_id/object", bucketHandler.GetObjectFromID)
			fileRoutes.GET("/", bucketHandler.GetFileByPath)
		}
	}

	return r
}
