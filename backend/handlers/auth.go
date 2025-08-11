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

	err = db.Transaction(func(tx *gorm.DB) error {

		// Create default role if it doesn't exist
		defaultRole := models.Role{
			Name: "user",
		}
		if err := tx.First(&defaultRole).Error; err != nil {
			log.Printf("Error creating default role: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign default role"})
			return err
		}

		newUser := models.User{
			Email:        req.Email,
			PasswordHash: hashedPassword,
			Username:     req.Username,
			IsActive:     true,
			RoleID:       defaultRole.ID,
		}
		if err := tx.Create(&newUser).Error; err != nil {
			log.Printf("Error saving new user to database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
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
		}
		if err := tx.Create(&session).Error; err != nil {
			log.Printf("Error creating session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return err
		}

		// Set cookies for session
		c.SetCookie("session", session.ID.String(), 3600, "/", "localhost", false, true)

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

	// Create a session for the user
	session := models.Session{
		UserID:    user.ID,
		Access:    accessTokenString,
		Refresh:   refreshTokenString,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
		IsValid:   true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := db.Create(&session).Error; err != nil {
		log.Printf("Error creating session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.SetCookie("session", session.ID.String(), 3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{"message": "User logged in successfully"})
}

func Logout(c *gin.Context) {
	db := config.GetDB()
	session, err := c.Cookie("session")
	if err != nil {
		log.Printf("Error retrieving session cookie: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session not found"})
		return
	}

	// Retrieve the session from the database
	var sessionModel models.Session
	if err := db.Find(&sessionModel, "id = ?", session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
			return
		}
		log.Printf("Database error retrieving session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	// Check if the session is already invalidated
	if !sessionModel.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session already invalidated"})
		return
	}

	c.SetCookie("session", "", -1, "/", "localhost", false, true)

	// Invalidate the session
	sessionModel.IsValid = false
	if err := db.Save(&sessionModel).Error; err != nil {
		log.Printf("Error invalidating session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate session"})
		return
	}

	c.JSON(200, gin.H{"message": "User logged out successfully"})
}
