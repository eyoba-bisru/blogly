package handlers

import (
	"net/http"

	"github.com/eyoba-bisru/blogly/config"
	"github.com/eyoba-bisru/blogly/helpers"
	"github.com/eyoba-bisru/blogly/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Validate the ID format
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var post models.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func GetPostBySlug(c *gin.Context) {
	db := config.GetDB()
	slug := c.Param("slug")

	var post models.Post
	if err := db.First(&post, "slug = ?", slug).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	db := config.GetDB()

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	post.UserID = c.MustGet("userID").(uuid.UUID)
	post.Slug = helpers.Slugify(post.Title)
	post.Published = false

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
	db := config.GetDB()
	id := c.Param("id")

	// Validate the ID format
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var post models.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	post.Slug = helpers.Slugify(post.Title)

	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	db := config.GetDB()
	id := c.Param("id")

	// Validate the ID format
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var post models.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func PublishPost(c *gin.Context) {
	db := config.GetDB()
	id := c.Param("id")

	userID := c.MustGet("userID").(uuid.UUID)

	// Check if the user has permission to publish posts
	if helpers.HasPermission(userID, "publish_posts") {

		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
			return
		}

		var post models.Post
		if err := db.First(&post, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		post.Published = true

		if err := db.Save(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post published successfully"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to publish posts"})
		return
	}

}
