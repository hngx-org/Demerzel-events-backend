package db

import "demerzel-events/internal/models"

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.UserGroup{},
		&models.Event{},
		&models.GroupEvent{},
		&models.InterestedEvent{},
	)
	return err
}
