package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	BaseModel
	LastChat     time.Time `gorm:"not null"`
	FirstCreated time.Time `gorm:"not null"`
	FileUploadID uuid.UUID
	FileUpload   FileUpload
	Messages     []Message
}
