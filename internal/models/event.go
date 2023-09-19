package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Creator     uuid.UUID `json:"creator"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	EventCreator User `gorm:"foreignKey:Creator"`
}

type InterestedEvent struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId  uuid.UUID `json:"user_id"`
	EventId uuid.UUID `json:"event_id"`

	User  User  `gorm:"foreignKey:User"`
	Event Event `gorm:"foreignKey:Event"`
}

type GroupEvent struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	GroupId uuid.UUID `json:"group_id"`
	EventId uuid.UUID `json:"event_id"`

	// needs Group model.
	Event Event `gorm:"foreignKey:Event"`
}
