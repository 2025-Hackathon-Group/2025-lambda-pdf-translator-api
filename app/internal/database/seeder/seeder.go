package seeder

import (
	"context"

	"gorm.io/gorm"
)

// Seeder defines the interface for all seeders
type Seeder interface {
	// Name returns the name of the seeder
	Name() string
	// Run executes the seeding logic
	Run(ctx context.Context, db *gorm.DB) error
}

// BaseSeeder provides common functionality for seeders
type BaseSeeder struct {
	DB *gorm.DB
}

// NewBaseSeeder creates a new base seeder
func NewBaseSeeder(db *gorm.DB) *BaseSeeder {
	return &BaseSeeder{
		DB: db,
	}
}
