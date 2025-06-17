package models

type Organization struct {
	BaseModel
	Name  string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Users []User
}
