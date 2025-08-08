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

	var count int64

	if posts.Count(&count); count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No posts found"})
		return
	}

	c.JSON(http.StatusOK, posts)

}

func GetPostByID(c *gin.Context) {
	db := config.GetDB()
	id := c.Param("id")

	var post models.Post
	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}
