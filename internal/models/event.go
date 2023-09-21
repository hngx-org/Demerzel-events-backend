package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewEvent struct {
	CreatorId   string `json:"creator" gorm:"type:varchar(255)"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Event struct {
	Id          string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	CreatorId   string `json:"creator_id" gorm:"type:varchar(255)"`
	Thumbnail   string `json:"thumbnail"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`

	Creator User `gorm:"foreignKey:CreatorId"`
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

func (group *Group) GetGroupEvent(tx *gorm.DB) (*[]Event, error){
	var events []Event
   
	err := tx.Table("group_events").Select("events.id, events.title, events.description, events.creator, events.location, events.start_date, events.end_date, events.start_time, events.end_time, events.created_at, events.updated_at").Joins("JOIN events on events.id = group_events.event_id").Where("group_events.group_id = ?", group.Id).Scan(&events).Error
    
	if err != nil {
		return nil, err
	}

	return &events, nil
}
func CreateEvent(tx *gorm.DB, event *NewEvent) (*Event, error) {

	request := Event{
		CreatorId:   event.CreatorId,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}

	err := tx.Model(Event{}).Create(&request)

	fmt.Print(event)

	if err.Error != nil {
		fmt.Print(err)

		return &Event{}, err.Error
	}
	return &request, nil
}

// ListAllEvents retrieves all events.
func ListEvents(tx *gorm.DB) ([]Event, error) {
	var events []Event

	err := tx.Order("start_date, start_time").Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}
