package main

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/database"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/router"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load environment variables: %v\n", err)
		return
	}

	// Initialize and connect to the database
	if err := database.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	// Initialize and start the server
	r := router.NewRouter(database.GetDB())

	fmt.Println("Server started on port " + os.Getenv("PORT"))
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
