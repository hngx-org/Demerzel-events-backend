package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	Id        string      `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name      string      `json:"name" validate:"required"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Members   []UserGroup `json:"members" gorm:"foreignkey:GroupID;association_foreignkey:ID"`
}

type UserGroup struct {
	Id        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255)"`
	GroupID   string    `json:"group_id" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	g.Id = uuid.NewString()

	return nil
}

func (uG *UserGroup) BeforeCreate(tx *gorm.DB) error {
	uG.Id = uuid.NewString()

	return nil
}

type Comment struct {
	Author    string    `json:"author" gorm:"primaryKey;type:varchar(500)"`
	Content   string    `json:"content" gorm:"type:text"`
	Post_id   string    `json:"post_id" gorm:"foreignkey:PostID;type:varchar(250)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
