package main

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/database"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load environment variables: %v\n", err)
		return
	}

	if err := database.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	if err := database.AutoMigrate(); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return
	}

	fmt.Println("Database migrated successfully")
}
