package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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
