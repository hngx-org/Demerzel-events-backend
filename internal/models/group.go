package models

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "time"
)

type Group struct {
    ID         string       `json:"id" gorm:"primaryKey;type:varchar(255)"`
    Name       string       `json:"name" validate:"required"`
    CreatedAt  time.Time    `json:"created_at"`
    UpdatedAt  time.Time    `json:"updated_at"`
    Members    []User      `gorm:"many2many:user_groups;"`
    GroupEvents []GroupEvent `gorm:"foreignKey:GroupID"`
}

type UserGroup struct {
    ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
    UserID    string    `json:"user_id" gorm:"type:varchar(255)"`
    GroupID   string    `json:"group_id" gorm:"type:varchar(255)"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    User      User      `gorm:"foreignKey:UserID"`
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
    g.ID = uuid.NewString()
    return nil
}

func (uG *UserGroup) BeforeCreate(tx *gorm.DB) error {
    uG.ID = uuid.NewString()
    return nil
}
