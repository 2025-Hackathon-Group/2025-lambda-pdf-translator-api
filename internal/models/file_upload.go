package models

import (
	"time"

	"github.com/google/uuid"
)

type FileUpload struct {
	BaseModel
	FileName        string    `gorm:"not null"`
	OriginalName    string    `gorm:"not null"`
	FileSize        int64     `gorm:"not null"`
	ContentType     string    `gorm:"not null"`
	S3Bucket        string    `gorm:"not null"`
	S3Key           string    `gorm:"not null;unique"`
	S3Region        string    `gorm:"not null"`
	UploadedAt      time.Time `gorm:"not null"`
	Path            string    `gorm:"not null;unique"`
	UserID          uuid.UUID
	User            User
	ProcessingState string `gorm:"type:string;not null;default:'pending'"` // pending, processing, completed, failed
	Error           string `gorm:"type:text"`
}
