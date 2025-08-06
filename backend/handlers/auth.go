package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eyoba-bisru/blogly/backend/config"
	"github.com/eyoba-bisru/blogly/backend/helpers"
	"github.com/eyoba-bisru/blogly/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	db := config.GetDB()

	// Bind the request body to the RegisterRequest struct
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var user models.User
	if err := db.First(&user, "email = ?", req.Email).Error; err == nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	} else if err.Error() != "record not found" {

		log.Printf("Database error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Println("SECRET_KEY environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}
	// Generate access tokens
	accessTokenString, err := helpers.AccessTokenGenerate(req.Email, secretKey)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	// Generate refresh token
	refreshTokenString, err := helpers.RefreshTokenGenerate(req.Email, secretKey)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie("access", accessTokenString, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh", refreshTokenString, 24*30*3600, "/", "localhost", false, true)

	err = db.Transaction(func(tx *gorm.DB) error {
		newUser := models.User{
			Email:        req.Email,
			PasswordHash: hashedPassword,
			Username:     req.Username,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := tx.Create(&newUser).Error; err != nil {
			log.Printf("Error saving new user to database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return err
		}

		// Create default role if it doesn't exist
		defaultRole := models.Role{
			Name: "user",
		}
		if err := tx.First(&defaultRole).Error; err != nil {
			log.Printf("Error creating default role: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign default role"})
			return err
		}

		// Assign default role to the user
		err = tx.Model(&newUser).Association("Roles").Append(&defaultRole)
		if err != nil {
			log.Printf("Error assigning role to user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
			return err
		}

		// Create a session for the user
		session := models.Session{
			UserID:    newUser.ID,
			Access:    accessTokenString,
			Refresh:   refreshTokenString,
			UserAgent: c.Request.UserAgent(),
			IPAddress: c.ClientIP(),
			IsValid:   true,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := tx.Create(&session).Error; err != nil {
			log.Printf("Error creating session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	// Login logic here
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.First(&user, "email = ?", req.Email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		log.Printf("Database error during login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	if !helpers.VerifyPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Println("SECRET_KEY environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}
	// Generate access tokens
	accessTokenString, err := helpers.AccessTokenGenerate(user.Email, secretKey)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	// Generate refresh token
	refreshTokenString, err := helpers.RefreshTokenGenerate(user.Email, secretKey)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie("access", accessTokenString, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh", refreshTokenString, 24*30*3600, "/", "localhost", false, true)
	// Create a session for the user
	session := models.Session{
		UserID:    user.ID,
		Access:    accessTokenString,
		Refresh:   refreshTokenString,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
		IsValid:   true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&session).Error; err != nil {
		log.Printf("Error creating session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(200, gin.H{"message": "User logged in successfully"})
}

func Logout(c *gin.Context) {
	// Logout logic here
	c.JSON(200, gin.H{"message": "User logged out successfully"})
}
