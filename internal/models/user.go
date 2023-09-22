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
	Events           []Event     `json:"-" gorm:"foreignKey:CreatorId"`
	InterestedEvents []Event     `json:"-" gorm:"many2many:interested_events;"`
	UserGroup        []UserGroup `json:"-" gorm:"foreignkey:UserID;association_foreignkey:ID"`
	Comments         []Comment   `json:"-" gorm:"foreignKey:CreatorId"`
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

func (u *User) GetUserByID(tx *gorm.DB) error {
	result := tx.First(&u, "id=?", u.Id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
