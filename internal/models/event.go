package models

import (
	"fmt"
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

type EventGroupReponse struct {
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
	GroupId     []string  `json:"group_id"`
}
type Event struct {
	Id          string       `json:"id" gorm:"primaryKey;type:varchar(255)"`
	CreatorId   string       `json:"creator_id" gorm:"type:varchar(255)"`
	Thumbnail   string       `json:"thumbnail"`
	Location    string       `json:"location"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	StartTime   string       `json:"start_time"`
	EndTime     string       `json:"end_time"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	Creator     *User        `json:"creator" gorm:"foreignKey:CreatorId"`
	Comments    []Comment    `json:"comments"`
	GroupEvents []GroupEvent `json:"group_events"`
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

func (group *Group) GetGroupEvent(tx *gorm.DB) (*[]Event, error) {
	var events []Event

	err := tx.Table("group_events").Select("events.id, events.title, events.description, events.creator, events.location, events.start_date, events.end_date, events.start_time, events.end_time, events.created_at, events.updated_at").Joins("JOIN events on events.id = group_events.event_id").Where("group_events.group_id = ?", group.ID).Scan(&events).Error

	if err != nil {
		return nil, err
	}

	return &events, nil
}

func CreateEvent(tx *gorm.DB, event *NewEvent) (*Event, error) {

	request := Event{
		CreatorId:   event.CreatorId,
		Thumbnail:   event.Thumbnail,
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
	fmt.Print(event)
	fmt.Print(event.GroupId)
	if len(event.GroupId) >= 1 {
		var groups []string
		for i := 0; i < len(event.GroupId); i++ {

			newGroupEvent := NewGroupEvent{
				EventId: request.Id,
				GroupId: event.GroupId[i],
			}
			res, err := CreateGroupEvent(tx, &newGroupEvent)

			if err != nil {
				fmt.Print(res)

				return nil, err
			}
			groups = append(groups, res.Id)
		}
		fmt.Print(groups)

		r, g, e := GetEventByParameter(tx, request.Id, "")
		fmt.Print(g)
		if e != nil {
			fmt.Print(r)

			return nil, e
		}

		return r, nil
	}

	return &request, nil
}

func CreateGroupEvent(tx *gorm.DB, groupEvent *NewGroupEvent) (*GroupEvent, error) {
	request := GroupEvent{
		EventId: groupEvent.EventId,
		GroupId: groupEvent.GroupId,
	}

	result := tx.Model(GroupEvent{}).Create(&request)

	if result.Error != nil {
		fmt.Print(result)

		return &GroupEvent{}, result.Error
	}

	return &request, nil
}

// retrieve an event using its ID
func GetEventByParameter(tx *gorm.DB, eventID string, date string) (*Event, []Event, error) {
	var event Event
	var events []Event

	if eventID != "" {

		err := tx.Where("id = ?", eventID).Preload("Creator").Preload("Comments").Preload("GroupEvents.Group").First(&event).Error

		if err != nil {
			return &Event{}, nil, err
		}

	}
	if date != "" {

		err := tx.Where("start_date = ?", date).Order("start_date, start_time").Preload("Creator").Preload("Comments").Preload("GroupEvents.Group").First(&events).Error

		if err != nil {
			return nil, []Event{}, err
		}

	}

	return &event, events, nil
}

// ListAllEvents retrieves all events.

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
