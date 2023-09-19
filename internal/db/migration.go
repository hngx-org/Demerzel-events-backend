package db

import "demerzel-events/internal/models"

func Migrate() {
	DB.AutoMigrate(&models.User{})
}
