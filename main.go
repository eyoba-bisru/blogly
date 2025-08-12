package main

import (
	"log"

	"github.com/eyoba-bisru/blogly/config"
	"github.com/eyoba-bisru/blogly/handlers"
	"github.com/eyoba-bisru/blogly/middlewares"
	"github.com/eyoba-bisru/blogly/seed"
	"github.com/gin-contrib/cors"
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

	// Seed roles and permissions
	seed.SeedRolesAndPermissions()

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	v1 := r.Group("/api/v1")
	{
		authV1 := v1.Group("/auth")
		{
			authV1.POST("/register", handlers.Register)
			authV1.POST("/login", handlers.Login)
			authV1.POST("/logout", handlers.Logout)
		}

		postsV1 := v1.Group("/posts")
		{
			postsV1.GET("/", handlers.GetPosts)
			postsV1.GET("/:id", handlers.GetPostByID)
			postsV1.GET("/slug/:slug", handlers.GetPostBySlug)
			postsV1.GET("/:id/comments", handlers.GetPostComments)
			postsV1.Use(middlewares.AuthMiddleware())
			postsV1.POST("/", handlers.CreatePost)
			postsV1.PATCH("/:id", handlers.UpdatePost)
			postsV1.DELETE("/:id", handlers.DeletePost)
			postsV1.POST("/publish/:id", handlers.PublishPost)
			postsV1.POST("/:id/comments", handlers.CreateComment)
		}

	}

	r.Run() // listen and serve on ":8080"
}
