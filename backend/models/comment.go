package models

import (
	"github.com/google/uuid"
)

type Comment struct {
	BaseModel
	Content string    `gorm:"not null" json:"content"`
	PostID  uuid.UUID `gorm:"not null" json:"post_id"`
	Post    Post      `gorm:"foreignKey:PostID" json:"post"`
	UserID  uuid.UUID `gorm:"not null" json:"user_id"`
	User    User      `gorm:"foreignKey:UserID" json:"user"`
}
