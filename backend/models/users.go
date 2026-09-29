package models

type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"` // "admin" or "operator"
}