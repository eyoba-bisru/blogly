package models

type Category struct {
	BaseModel
	Name  string `gorm:"not null;unique" json:"name"`
	Slug  string `gorm:"not null;unique" json:"slug"`
	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts"`
}
