package models

type Role struct {
	BaseModel
	Name        string       `gorm:"not null;unique" json:"name"`
	Description *string      `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}
