package config

import (
	"os"

	"github.com/eyoba-bisru/blogly/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() *gorm.DB {

	dns := os.Getenv("DATABASE_URL")
	if dns == "" {
		dns = "host=localhost user=postgres password=yourpassword dbname=blogly port=5432 sslmode=disable"
	}

	DB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = DB

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(models.Post{})
	DB.AutoMigrate(&models.Comment{})
	DB.AutoMigrate(&models.Role{})
	DB.AutoMigrate(&models.Permission{})
	DB.AutoMigrate(&models.Session{})
	return DB
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("Database connection is not initialized")
	}
	return db
}
