package models

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewComment struct {
	Body    string   `json:"body" validate:"required,min=1"`
	EventId string   `json:"event_id"`
	Images  []string `json:"images"`
}

type Images []string
type Comment struct {
	Id        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Body      string    `json:"body" validate:"required,min=1"`
	EventId   string    `json:"event_id"`
	Images    Images    `json:"images" gorm:"type:TEXT"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Creator   *User     `json:"creator" gorm:"foreignKey:CreatorId"`
	CreatorId string    `json:"creator_id" gorm:"type:varchar(255)"`
}

type CommentResponse struct {
	Id        string         `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Body      string         `json:"body" validate:"required,min=1"`
	EventId   string         `json:"event_id"`
	Images    Images         `json:"images" gorm:"type:TEXT"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Creator   CommentCreator `json:"creator"`
}

type CommentCreator struct {
	Id     string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Avatar string `gorm:"column:avatar" json:"avatar"`
}

type UpdateComment struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	c.Id = uuid.NewString()
	return nil
}

func (o *Images) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(string(bytes), ",")
	return nil
}

func (o Images) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}
