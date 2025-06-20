package seeder

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"context"

	"gorm.io/gorm"
)

type OrganisationSeeder struct {
	*BaseSeeder
}

func NewOrganisationSeeder(db *gorm.DB) *OrganisationSeeder {
	return &OrganisationSeeder{
		BaseSeeder: NewBaseSeeder(db),
	}
}

func (s *OrganisationSeeder) Name() string {
	return "organisation_seeder"
}

func (s *OrganisationSeeder) Run(ctx context.Context, db *gorm.DB) error {
	organisations := []models.Organisation{
		{
			Name:  "Default Organisation",
			Email: "admin@default.org",
		},
	}

	for _, org := range organisations {
		if err := db.FirstOrCreate(&org, models.Organisation{Email: org.Email}).Error; err != nil {
			return err
		}
	}

	return nil
}
