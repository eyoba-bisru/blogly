package seed

import (
	"github.com/eyoba-bisru/blogly/backend/config"
	"github.com/eyoba-bisru/blogly/backend/models"
)

func SeedRolesAndPermissions() {

	db := config.GetDB()

	// Permissions
	perm := models.Permission{Name: "read_posts"}
	db.FirstOrCreate(&perm, perm)

	// Role
	userRole := models.Role{Name: "user"}
	db.FirstOrCreate(&userRole, userRole)
	db.Model(&userRole).Association("Permissions").Append(&perm)

	// Category
	category := models.Category{Name: "General", Slug: "general"}
	db.FirstOrCreate(&category, category)

	category2 := models.Category{Name: "Technology", Slug: "technology"}
	db.FirstOrCreate(&category2, category2)
}
