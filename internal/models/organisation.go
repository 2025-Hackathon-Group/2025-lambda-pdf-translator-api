package models

type Organisation struct {
	BaseModel
	Name  string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Users []User
}
