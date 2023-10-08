package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID        string      `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name      string      `json:"name" validate:"required"`
	Image     string      `json:"image"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Members   []UserGroup `json:"members" gorm:"foreignkey:GroupID;association_foreignkey:ID;constraint:OnDelete:CASCADE"`
	Events    []Event     `json:"events,omitempty" gorm:"many2many:group_events"`
}

type GroupResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Image     string          `json:"image"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Members   []UserGroup     `json:"members"`
	Events    []EventResponse `json:"events"`
}

type UserGroup struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255)"`
	GroupID   string    `json:"group_id" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user" gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

type UpdateGroupRequest struct {
	Name string `json:"name"`
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	g.ID = uuid.NewString()

	return nil
}

func (uG *UserGroup) BeforeCreate(tx *gorm.DB) error {
	uG.ID = uuid.NewString()

	return nil
}
