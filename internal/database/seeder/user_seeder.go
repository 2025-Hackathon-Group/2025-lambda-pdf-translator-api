package seeder

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/service"
	"context"
	"os"

	"gorm.io/gorm"
)

type UserSeeder struct {
	*BaseSeeder
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{
		BaseSeeder: NewBaseSeeder(db),
	}
}

func (s *UserSeeder) Name() string {
	return "user_seeder"
}

func (s *UserSeeder) Run(ctx context.Context, db *gorm.DB) error {
	userService := service.NewUserService(os.Getenv("JWT_SECRET"))
	hashedPassword, err := userService.HashPassword("password")
	if err != nil {
		return err
	}

	org := models.Organisation{}
	db.First(&org)

	users := []models.User{
		{
			Name:           "Default User",
			Email:          "saml@everbit.dev",
			Password:       hashedPassword,
			OrganisationID: org.ID,
		},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			return err
		}
	}

	return nil
}
