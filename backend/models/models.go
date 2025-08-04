package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"password_hash"`
	CreatedAt    int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

type Category struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null;unique" json:"name"`
	Slug      string `gorm:"not null;unique" json:"slug"`
	Posts     []Post `gorm:"foreignKey:CategoryID" json:"posts"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

type Post struct {
	gorm.Model
	ID         uint     `gorm:"primaryKey" json:"id"`
	Title      string   `gorm:"not null" json:"title"`
	Content    string   `gorm:"not null" json:"content"`
	Slug       string   `gorm:"not null;unique" json:"slug"`
	UserID     uint     `gorm:"not null" json:"user_id"`
	User       User     `gorm:"foreignKey:UserID" json:"user"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
	Published  bool     `gorm:"default:false" json:"published"`
	CreatedAt  int64    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64    `json:"updated_at" gorm:"autoUpdateTime"`
}

type Comment struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	Content   string `gorm:"not null" json:"content"`
	PostID    uint   `gorm:"not null" json:"post_id"`
	Post      Post   `gorm:"foreignKey:PostID" json:"post"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}
