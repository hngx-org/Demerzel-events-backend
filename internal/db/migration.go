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
		&models.Comment{},
		&models.Notification{},
		&models.UserNotification{},
		&models.Reaction{},
		&models.Tag{},
		&models.GroupTag{},
		&models.NotificationSetting{},
	)
	return err
}
