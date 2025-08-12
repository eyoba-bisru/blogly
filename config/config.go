package config

import (
	"os"

	"github.com/eyoba-bisru/blogly/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {

	dns := os.Getenv("DATABASE_URL")
	if dns == "" {
		dns = "host=localhost user=postgres password=yourpassword DBname=blogly port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db

	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

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
	if DB == nil {
		panic("Database connection is not initialized")
	}
	return DB
}
