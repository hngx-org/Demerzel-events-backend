package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id     string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Email  string `gorm:"column:email;unique" json:"email"`
	Avatar string `gorm:"column:avatar" json:"avatar"`

	// Events Relationship
	Events           []Event     `gorm:"foreignKey:Creator"`
	InterestedEvents []Event     `gorm:"many2many:interested_events;"`
	UserGroup        []UserGroup `json:"user_group" gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

func NewUser(name string, email string, avatar string) *User {
	return &User{
		Name:   name,
		Email:  email,
		Avatar: avatar,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.NewString()

	return nil
}
