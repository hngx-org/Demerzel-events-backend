package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewEvent struct {
	CreatorId   string `json:"creator" gorm:"type:varchar(255)"`
	Thumbnail   string `json:"thumbnail"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
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
	Creator     *User     `json:"creator" gorm:"foreignKey:CreatorId"`
	Attendees   []User    `json:"attendees" gorm:"many2many:interested_events"`
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

	result := tx.Model(Event{}).Create(&request)

	if result.Error != nil {
		fmt.Print(result)

		return &Event{}, result.Error
	}

	return &request, nil
}

// GetEventByID retrieve an event using its ID
func GetEventByID(tx *gorm.DB, eventID string) (*Event, error) {
	var event Event

	err := tx.Where("id = ?", eventID).Preload("Creator").Preload("Attendees").First(&event).Error

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func AttachUserToEvent(tx *gorm.DB, userID string, eventID string) (*Event, error) {
	err := tx.Model(&InterestedEvent{}).FirstOrCreate(&InterestedEvent{UserId: userID, EventId: eventID}).Error
	if err != nil {
		return nil, err
	}

	return GetEventByID(tx, eventID)
}

// ListEvents retrieves all events.
func ListEvents(tx *gorm.DB) ([]Event, error) {
	var events []Event

	err := tx.Order("start_date, start_time").Preload("Creator").Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil

}

func ListEventsInGroups(tx *gorm.DB, groupIDs []string) ([]Event, error) {
	var events []Event

	err := tx.Model(&Event{}).
		Joins("JOIN group_events ON events.id = group_events.event_id").
		Where("group_events.group_id IN ?", groupIDs).
		Distinct("events.id").
		Select("events.id, events.creator_id, events.thumbnail," +
			"events.location, events.title, events.description, " +
			"events.start_date,events.end_date,events.start_time,events.end_time,events.created_at, events.updated_at").
		Preload("Creator").Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}
