package handlers

import (
	"net/http"

	"github.com/eyoba-bisru/blogly/backend/config"
	"github.com/eyoba-bisru/blogly/backend/models"
	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	db := config.GetDB()

	posts := db.Find(&models.Post{})
	if posts.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, posts)

}
