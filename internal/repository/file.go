package repository

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"

	"gorm.io/gorm"
)

type FileRepository interface {
	CreateFile(file *models.FileUpload) error
	GetFile(id string) (*models.FileUpload, error)
	GetFileByPath(path string) (*models.FileUpload, error)
}

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

// CreateFile creates a new file
func (r *fileRepository) CreateFile(file *models.FileUpload) error {
	return r.db.Create(file).Error
}

// GetFile returns a file by id
func (r *fileRepository) GetFile(id string) (*models.FileUpload, error) {
	var file models.FileUpload
	if err := r.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFileByPath returns a file by path
func (r *fileRepository) GetFileByPath(path string) (*models.FileUpload, error) {
	var file models.FileUpload
	if err := r.db.First(&file, "path = ?", path).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
