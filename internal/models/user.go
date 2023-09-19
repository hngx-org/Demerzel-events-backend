package models

import (
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
}
