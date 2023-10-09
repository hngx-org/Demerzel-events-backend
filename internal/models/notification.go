package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Type      string    `json:"type" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.ID = uuid.NewString()

	return nil
}

type UserNotification struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID         string       `json:"user_id" gorm:"type:varchar(255)"`
	NotificationID string       `json:"notification_id" gorm:"type:varchar(255)"`
	Read           bool         `json:"read" gorm:"default:false"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Notification   Notification `json:"notification" gorm:"foreignkey:NotificationID;association_foreignkey:ID"`
	User           User         `json:"user" gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

func (uN *UserNotification) BeforeCreate(tx *gorm.DB) error {
	uN.ID = uuid.NewString()

	return nil
}

type NotificationSetting struct {
	ID       string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID   string `json:"user_id" gorm:"type:varchar(255)"`
	User     User   `json:"user" gorm:"foreignKey:UserID;association_foreignkey:ID"`
	Email    bool   `json:"email" gorm:"default:true"`
	Event    bool   `json:"event" gorm:"default:true"`
	Group    bool   `json:"group" gorm:"default:true"`
	Reminder bool   `json:"reminder" gorm:"default:true"`
}

func (nS *NotificationSetting) BeforeCreate(tx *gorm.DB) error {
	nS.ID = uuid.NewString()

	return nil
}
