package router

import (
	"github.com/gin-gonic/gin"
)

// NewRouter initializes and returns a new Gin router
func NewRouter() *gin.Engine {
	r := gin.Default()

	// Example route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}
