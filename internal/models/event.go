package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewEvent struct {
	CreatorId   string   `json:"creator" gorm:"type:varchar(255)"`
	Thumbnail   string   `json:"thumbnail"`
	GroupId     []string `json:"group_id"`
	Location    string   `json:"location"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
}

type NewGroupEvent struct {
	EventId string `json:"event_id"`
	GroupId string `json:"group_id"`
}

type Event struct {
	Id          string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	CreatorId   string    `json:"creator_id" gorm:"type:varchar(255)"`
	Thumbnail   string    `json:"thumbnail"`
	Location    string    `json:"location"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Creator     *User     `json:"creator,omitempty" gorm:"foreignKey:CreatorId"`
	Attendees   []User    `json:"attendees,omitempty" gorm:"many2many:interested_events"`
	Groups      []Group   `json:"groups,omitempty" gorm:"many2many:group_events"`
	Comments    []Comment `json:"comments,omitempty" gorm:"foreignKey:EventId;constraint:OnDelete:CASCADE"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.Id = uuid.NewString()

	return nil
}

type InterestedEvent struct {
	Id      string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserId  string `json:"user_id" gorm:"type:varchar(255)"`
	EventId string `json:"event_id" gorm:"type:varchar(255)"`

	User  User  `json:"user" gorm:"foreignKey:UserId"`
	Event Event `json:"event" gorm:"foreignKey:EventId"`
}

func (iE *InterestedEvent) BeforeCreate(tx *gorm.DB) error {
	iE.Id = uuid.NewString()

	return nil
}

type GroupEvent struct {
	gorm.Model
	Id      string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	GroupId string `json:"group_id" gorm:"type:varchar(255)"`
	EventId string `json:"event_id" gorm:"type:varchar(255)"`

	Group Group `gorm:"foreignKey:GroupId; constraint:OnUpdate:CASCADE"`
	Event Event `gorm:"foreignKey:EventId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (gE *GroupEvent) BeforeCreate(tx *gorm.DB) error {
	gE.Id = uuid.NewString()

	return nil
}
