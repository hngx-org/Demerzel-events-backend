package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	Id        string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
	EventId   string    `json:"event_id"`
	Images    []Image   `json:"images" gorm:"foreignkey:CommentId"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Image struct {
	Id        string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	CommentId string `json:"comment_id"`
	ImageUrl  string `json:"image_url"`
}

type UpdateComment struct {
	EventId string `json:"event_id"`
	Body    string `json:"body"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	c.Id = uuid.NewString()
	return nil
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	i.Id = uuid.NewString()
	return nil
}
