package models

type Permission struct {
	BaseModel
	Name        string  `gorm:"not null;unique" json:"name"`
	Description *string `json:"description"`
	Roles       []Role  `gorm:"many2many:role_permissions" json:"roles"`
}
