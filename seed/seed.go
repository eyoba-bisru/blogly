package seed

import (
	"github.com/eyoba-bisru/blogly/config"
	"github.com/eyoba-bisru/blogly/helpers"
	"github.com/eyoba-bisru/blogly/models"
)

func SeedRolesAndPermissions() {

	db := config.GetDB()

	// Permissions
	perm := models.Permission{Name: "read_posts"}
	db.FirstOrCreate(&perm, perm)

	perm2 := models.Permission{Name: "write_posts"}
	db.FirstOrCreate(&perm2, perm2)

	perm3 := models.Permission{Name: "delete_posts"}
	db.FirstOrCreate(&perm3, perm3)

	perm4 := models.Permission{Name: "publish_posts"}
	db.FirstOrCreate(&perm4, perm4)

	// Role
	userRole := models.Role{Name: "user"}
	db.FirstOrCreate(&userRole, userRole)
	db.Model(&userRole).Association("Permissions").Append([]models.Permission{perm, perm2, perm3})

	adminRole := models.Role{Name: "admin"}
	db.FirstOrCreate(&adminRole, adminRole)
	db.Model(&adminRole).Association("Permissions").Append([]models.Permission{perm, perm2, perm3, perm4})

	// Category
	category := models.Category{Name: "General", Slug: helpers.Slugify("general")}
	db.FirstOrCreate(&category, category)

	category2 := models.Category{Name: "Technology", Slug: helpers.Slugify("technology")}
	db.FirstOrCreate(&category2, category2)
}
