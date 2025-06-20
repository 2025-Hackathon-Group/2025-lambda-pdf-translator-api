package main

import (
	"context"
	"log"
	"os"

	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/database"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/database/seeder"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seed",
	Short: "Database seeding tool",
	Long:  `A tool for seeding the database with initial data`,
}

func init() {
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(resetCmd)
}

var seedCmd = &cobra.Command{
	Use:   "run",
	Short: "Run database seeding",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Connect to database
		if err := database.Connect(); err != nil {
			return err
		}

		// Run migrations
		if err := database.AutoMigrate(); err != nil {
			return err
		}

		// Run seeders
		if err := runSeeders(ctx); err != nil {
			return err
		}

		return nil
	},
}

func runSeeders(ctx context.Context) error {
	db := database.GetDB()

	// Initialize seeders
	seeders := []seeder.Seeder{
		seeder.NewOrganisationSeeder(db),
		seeder.NewUserSeeder(db),
	}

	// Run each seeder
	for _, s := range seeders {
		log.Printf("Running seeder: %s", s.Name())
		if err := s.Run(ctx, db); err != nil {
			return err
		}
		log.Printf("Completed seeder: %s", s.Name())
	}

	return nil
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Connect to database
		if err := database.Connect(); err != nil {
			return err
		}

		// Run migrations
		if err := database.AutoMigrate(); err != nil {
			return err
		}

		return resetDatabase(ctx)
	},
}

func resetDatabase(ctx context.Context) error {

	db := database.GetDB()
	return db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
