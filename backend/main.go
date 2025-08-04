package main

import (
	"log"

	"github.com/eyoba-bisru/blogly/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	r.Run() // listen and serve on ":8080"
}
