package config

import (
	"os"

	"github.com/eyoba-bisru/blogly/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {

	dns := os.Getenv("DATABASE_URL")
	if dns == "" {
		dns = "host=localhost user=postgres password=yourpassword dbname=blogly port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Comment{})
	return db
}
