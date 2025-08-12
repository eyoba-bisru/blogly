package models

import (
	"github.com/google/uuid"
)

type Post struct {
	BaseModel
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"not null" json:"content"`
	Slug       string    `gorm:"not null;unique" json:"slug"`
	UserID     uuid.UUID `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	CategoryID uuid.UUID `gorm:"not null" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Published  bool      `gorm:"default:false" json:"published"`
}
