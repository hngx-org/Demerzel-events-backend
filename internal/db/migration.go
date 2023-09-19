package db

import "demerzel-events/internal/models"

func Migrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Comment{})
	DB.AutoMigrate(&models.Image{})
	DB.AutoMigrate(&models.Group{})
	DB.AutoMigrate(&models.UserGroup{})
}
