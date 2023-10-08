package types

import "demerzel-events/internal/models"

type NotificationTypes int

const (
	EventNotification NotificationTypes = iota
	SubscriptionNotification
	GroupNotification
)

type UserNotificationResponse struct {
	ID             string              `json:"id"`
	UserID         string              `json:"user_id"`
	NotificationID string              `json:"notification_id"`
	Read           bool                `json:"read"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
	Notification   models.Notification `json:"notification" gorm:"foreignkey:NotificationID;association_foreignkey:ID"`
}

type NotificationSettingsResponse struct {
	ID       string `json:"id"`
	Email    bool `json:"email"`
	Event    bool `json:"event"`
	Group    bool `json:"group"`
	Reminder bool `json:"reminder"`
}
