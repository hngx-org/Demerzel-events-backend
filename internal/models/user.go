package models

type User struct {
	// Add user fields
	UserGroup []UserGroup `json:"user_group" gorm:"foreignkey:UserID;association_foreignkey:ID"`
}
