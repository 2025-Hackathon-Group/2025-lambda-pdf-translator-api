package seeder

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"context"

	"gorm.io/gorm"
)

type OrganizationSeeder struct {
	*BaseSeeder
}

func NewOrganizationSeeder(db *gorm.DB) *OrganizationSeeder {
	return &OrganizationSeeder{
		BaseSeeder: NewBaseSeeder(db),
	}
}

func (s *OrganizationSeeder) Name() string {
	return "organization_seeder"
}

func (s *OrganizationSeeder) Run(ctx context.Context, db *gorm.DB) error {
	organizations := []models.Organization{
		{
			Name:  "Default Organization",
			Email: "admin@default.org",
		},
	}

	for _, org := range organizations {
		if err := db.FirstOrCreate(&org, models.Organization{Email: org.Email}).Error; err != nil {
			return err
		}
	}

	return nil
}
