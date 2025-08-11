package helpers

import (
	"log"

	"github.com/eyoba-bisru/blogly/backend/config"
	"github.com/eyoba-bisru/blogly/backend/models"
	"github.com/google/uuid"
)

func HasPermission(userID uuid.UUID, permission string) bool {
	db := config.GetDB()
	var user models.User

	var count int64

	err := db.Model(&user).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Joins("INNER JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("INNER JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("users.id = ? AND permissions.name = ?", userID, permission).
		Count(&count).Error

	if err != nil {
		log.Printf("Error checking permission for user %s: %v", userID, err)
		return false
	}

	return count > 0

}
