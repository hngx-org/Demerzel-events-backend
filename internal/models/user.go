package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Username string `json:"username" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"-" gorm:"not null"`


	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`


	Events           []Event     `gorm:"foreignKey:Creator"`
	InterestedEvents []Event     `gorm:"many2many:interested_events;"`
	UserGroups       []UserGroup `json:"user_groups" gorm:"many2many:user_groups;"`

}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.NewString()
	return nil
}

