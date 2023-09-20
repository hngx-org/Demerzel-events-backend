package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	Id          string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Creator     string    `json:"creator" gorm:"type:varchar(255)"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	EventCreator User `gorm:"foreignKey:Creator"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.Id = uuid.NewString()

	return nil
}

type InterestedEvent struct {
	gorm.Model
	Id      string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserId  string `json:"user_id" gorm:"type:varchar(255)"`
	EventId string `json:"event_id" gorm:"type:varchar(255)"`

	User  User  `gorm:"foreignKey:UserId"`
	Event Event `gorm:"foreignKey:EventId"`
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

	Group Group `gorm:"foreignKey:GroupId"`
	Event Event `gorm:"foreignKey:EventId"`
}

func (gE *GroupEvent) BeforeCreate(tx *gorm.DB) error {
	gE.Id = uuid.NewString()

	return nil
}

func (group *Group) GetGroupEvent(tx *gorm.DB) (*[]Event){
	var events []Event
   
	tx.Table("group_events").Select("events.id, events.title, events.description, events.creator, events.location, events.start_date, events.end_date, events.start_time, events.end_time, events.created_at, events.updated_at").Joins("JOIN events on events.id = group_events.event_id").Where("group_events.event_id = ?", group.Id).Scan(&events)
    
	return &events
}