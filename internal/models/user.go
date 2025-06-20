package models

import "github.com/google/uuid"

type User struct {
	BaseModel
	Email          string `gorm:"unique;not null"`
	Name           string `gorm:"not null"`
	Password       string `gorm:"not null"`
	ProfilePicture string
	OrganisationID uuid.UUID
	Organisation   Organisation
	FileUploads    []FileUpload
}
