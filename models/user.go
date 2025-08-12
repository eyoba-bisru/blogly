package models

import "github.com/google/uuid"

type User struct {
	BaseModel
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	RoleID       uuid.UUID `gorm:"not null" json:"role_id"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"role"`
	Sessions     []Session `gorm:"foreignKey:UserID" json:"sessions"`
	Posts        []Post    `gorm:"foreignKey:UserID" json:"posts"`
	Comments     []Comment `gorm:"foreignKey:UserID" json:"comments"`
}
