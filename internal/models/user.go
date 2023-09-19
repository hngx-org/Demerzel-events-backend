package models

import (
<<<<<<< HEAD
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    Id string `json:"id" gorm:"primaryKey;type:varchar(255)"`
    // Add user fields

    Name         string `gorm:"column:name" json:"name"`
    Email        string `gorm:"column:email;unique" json:"email"`
    AccessToken  string `gorm:"column:access_token" json:"access_token"`
    RefreshToken string `gorm:"column:refresh_token" json:"refresh_token"`
    Avatar       string `gorm:"column:avatar" json:"avatar"`
    // Events Relationship
    Events           []Event     `gorm:"foreignKey:Creator"`
    InterestedEvents []Event     `gorm:"many2many:interested_events;"`
    UserGroup        []UserGroup `json:"user_group" gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.Id = uuid.NewString()

    return nil
}

type UserResponse struct {
    Name   string `json:"name"`
    Email  string `json:"email"`
    Avatar string `json:"avatar"`
=======
	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;unique"`
	Name         string    `gorm:"type:varchar(255)"`
	Email        string    `gorm:"type:varchar(255);unique"`
	Password     string    `gorm:"type:varchar(255)"`
	AccessToken  string    `gorm:"type:varchar(255)"`
	RefreshToken string    `gorm:"type:varchar(255)"`
	Avatar       string    `gorm:"type:varchar(255)"`
}

type UserResponse struct {
	Username  string
	Email     string
	Token     string
	TokenType string
}

type UserSignUp struct {
	Name     string `gorm:"type:varchar(255)"`
	Email    string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
	Avatar   string `gorm:"type:varchar(255)"`
}

type UserLogin struct {
	Email    string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
}

type PasswordReset struct {
	Email           string `gorm:"type:varchar(255)"`
	CurrentPassword string `gorm:"type:varchar(255)"`
	NewPassword     string `gorm:"type:varchar(255)"`
	ConfirmPassword string `gorm:"type:varchar(255)"`
}

type UpdateUser struct {
	Email  string `gorm:"type:varchar(255)"`
	Name   string `gorm:"type:varchar(255)"`
	Avatar string `gorm:"type:varchar(255)"`
}

type ForgotPassword struct {
	Email string `gorm:"type:varchar(255)"`
>>>>>>> 53b284fe009603a470b7c03ac5f9ed99a31ac799
}
