package models

import (
	"time"
)

type Comment struct {
	Id        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
	EventId   string    `json:"event_id"`
	Images    []Image   `json:"images" gorm:"foreignkey:CommentId"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Image struct {
	Id        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	CommentId string    `json:"comment_id"`
	ImageUrl  string    `json:"image_url"`
}
