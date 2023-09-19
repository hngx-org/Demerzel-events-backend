package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	// Add user fields

	// Events Relationship
	Events           []Event     `gorm:"foreignKey:Creator"`
	InterestedEvents []Event     `gorm:"many2many:interested_events;"`
	UserGroup        []UserGroup `json:"user_group" gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.NewString()

	return nil
}
