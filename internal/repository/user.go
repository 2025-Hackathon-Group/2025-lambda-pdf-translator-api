package repository

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user with a hashed password
func (r *UserRepository) Create(name, email, password string) (*models.User, error) {
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password, // Password should be hashed before calling this method
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user's information
func (r *UserRepository) Update(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// GetOrganisations retrieves all organisations a user belongs to
func (r *UserRepository) GetOrganisations(userID uuid.UUID) ([]models.Organisation, error) {
	var orgs []models.Organisation
	err := r.db.Joins("JOIN users ON users.organisation_id = organisations.id").
		Where("users.id = ?", userID).
		Find(&orgs).Error
	return orgs, err
}

func (r *UserRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}
