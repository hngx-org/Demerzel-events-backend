package models

import (
	"time"
	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Members   []UserGroup `json:"members" gorm:"foreignkey:GroupID;association_foreignkey:ID"`
}

type UserGroup struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id"`
	GroupID   uuid.UUID `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
