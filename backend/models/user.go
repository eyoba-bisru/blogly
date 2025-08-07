package models

type User struct {
	BaseModel
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Roles        []Role    `gorm:"many2many:user_roles" json:"roles"`
	Sessions     []Session `gorm:"foreignKey:UserID" json:"sessions"`
	Posts        []Post    `gorm:"foreignKey:UserID" json:"posts"`
	Comments     []Comment `gorm:"foreignKey:UserID" json:"comments"`
}
