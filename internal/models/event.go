package models

import (
	"errors"
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
	Groups      []Group   `json:"groups" gorm:"many2many:group_events"`
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

	Group Group `gorm:"foreignKey:GroupId"`
	Event Event `gorm:"foreignKey:EventId"`
}

func (gE *GroupEvent) BeforeCreate(tx *gorm.DB) error {
	gE.Id = uuid.NewString()

	return nil
}

//
//func (g *Group) GetGroupEvents(tx *gorm.DB) (*[]Event, error) {
//	result := tx.Preload("Events").Where("id = ?", g.ID).First(g)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//
//	return &g.Events, nil
//}

func CreateEvent(tx *gorm.DB, event *NewEvent) (*Event, error) {

	createdEvent := Event{
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

	err := tx.Transaction(func(tx *gorm.DB) error {

		result := tx.Model(Event{}).Create(&createdEvent)

		if result.Error != nil {
			fmt.Print(result)

			return result.Error
		}

		if len(event.GroupId) >= 1 {
			var groups []string
			for i := 0; i < len(event.GroupId); i++ {

				newGroupEvent := NewGroupEvent{
					EventId: createdEvent.Id,
					GroupId: event.GroupId[i],
				}
				res, err := CreateGroupEvent(tx, &newGroupEvent)

				if err != nil {
					fmt.Print(res)

					return err
				}
				groups = append(groups, res.Id)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return GetEventByID(tx, createdEvent.Id)
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

// GetEventByID retrieve an event using its ID
func GetEventByID(tx *gorm.DB, eventID string) (*Event, error) {
	var event Event

	err := tx.Where("id = ?", eventID).Preload("Creator").Preload("Attendees").Preload("Groups").First(&event).Error

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

func DetachUserFromEvent(tx *gorm.DB, userID, eventID string) (*Event, error) {
	err := tx.Model(&InterestedEvent{}).Unscoped().Where(&InterestedEvent{UserId: userID, EventId: eventID}).Delete(&InterestedEvent{}).Error

	if err != nil {
		return nil, err
	}

	return GetEventByID(tx, eventID)
}

// ListEvents retrieves all events.
func ListEvents(tx *gorm.DB) ([]Event, error) {
	var events []Event

	err := tx.Order("start_date, start_time").Preload("Creator").Preload("Attendees").Preload("Groups").Find(&events).Error

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

func SubscribeUserToEvent(tx *gorm.DB, userID, eventID string) (*InterestedEvent, error) {
	var interestedEvent InterestedEvent
	result := tx.Where("event_id = ?", eventID).Where("user_id = ?", userID).First(&interestedEvent)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		interestedEvent = InterestedEvent{
			UserId:  userID,
			EventId: eventID,
		}

		result = tx.Create(&interestedEvent)
		if result.Error != nil {
			return nil, result.Error
		}

		return &interestedEvent, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return nil, fmt.Errorf("user already subscribed to event")
}

func UnsubscribeUserFromEvent(tx *gorm.DB, userID, eventID string) error {
	var interestedEvent InterestedEvent
	result := tx.Where("event_id = ?", eventID).Where("user_id = ?", userID).First(&interestedEvent)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not subscribed to event")
		}
		return result.Error
	}

	// Delete the UserGroup
	result = tx.Delete(&interestedEvent)
	return result.Error
}

func GetUserEventSubscriptions(tx *gorm.DB, userID string) (*[]Event, error) {
	var user User
	result := tx.Where("id = ?", userID).Preload("InterestedEvents").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user.InterestedEvents, nil
}
