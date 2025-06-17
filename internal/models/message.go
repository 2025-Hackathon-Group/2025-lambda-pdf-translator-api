package models

import (
	"time"

	"github.com/google/uuid"
)

type SenderType string

const (
	UserSender SenderType = "USER"
	AISender   SenderType = "AI"
)

type Message struct {
	BaseModel
	Sender  SenderType `gorm:"type:string;not null"`
	Content string     `gorm:"not null"`
	SentAt  time.Time  `gorm:"not null"`
	ChatID  uuid.UUID
	Chat    Chat
}
