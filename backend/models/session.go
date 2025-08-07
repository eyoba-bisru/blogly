package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	BaseModel
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	UserAgent string    `gorm:"not null" json:"user_agent"`
	IPAddress string    `gorm:"not null" json:"ip_address"`
	Access    string    `gorm:"not null;unique" json:"access"`
	Refresh   string    `gorm:"not null;unique" json:"refresh"`
	IsValid   bool      `gorm:"default:true" json:"is_valid"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}
