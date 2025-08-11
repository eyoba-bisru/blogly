package handlers

import (
	"net/http"

	"github.com/eyoba-bisru/blogly/backend/config"
	"github.com/eyoba-bisru/blogly/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPostComments(c *gin.Context) {
	db := config.GetDB()
	postID := c.Param("id")

	// Validate the post ID format
	if _, err := uuid.Parse(postID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var comments []models.Comment
	if err := db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	if len(comments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No comments found for this post"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
	db := config.GetDB()
	postID := c.Param("id")

	// Validate the post ID format
	if _, err := uuid.Parse(postID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	comment.PostID = uuid.MustParse(postID)
	comment.UserID = c.MustGet("userID").(uuid.UUID)

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}
