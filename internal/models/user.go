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



   
    // AccessToken  string `gorm:"column:access_token" json:"access_token"`
	// RefreshToken string `gorm:"column:refresh_token" json:"refresh_token"` 
    // Events Relationship
    

	     
	// Events Relationship

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


type UserResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type UserUpdate struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserData struct {
    Id string `json:"id" gorm:"primaryKey;type:varchar(255)"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Avatar string `json:"avatar"`
}
type UpdateUserStruct struct {
    Name   string `json:"name"`
    Email  string `json:"email"`
}

