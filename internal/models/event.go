package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Creator     string    `json:"creator" gorm:"type:varchar(255)"`
	Location    string    `json:"location"`
	OrganizerID uint      `json:"organizer_id"`
	ImageURL    string    `json:"image_url"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	EventCreator User `gorm:"foreignKey:Creator"`
	Groups       []Group `gorm:"many2many:group_events;"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.NewString()
	return nil
}

type InterestedEvent struct {
	gorm.Model
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID      string    `json:"user_id" gorm:"type:varchar(255)"`
	EventID     string    `json:"event_id" gorm:"type:varchar(255)"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	OrganizerID uint      `json:"organizer_id"`
	StartDate   time.Time `json:"start_date"`

	User  User  `gorm:"foreignKey:UserID"`
	Event Event `gorm:"foreignKey:EventID"`
}

func (iE *InterestedEvent) BeforeCreate(tx *gorm.DB) error {
	iE.ID = uuid.NewString()
	return nil
}

type GroupEvent struct {
	gorm.Model
	ID       string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	GroupID  string    `json:"group_id" gorm:"type:varchar(255)"`
	EventID  string    `json:"event_id" gorm:"type:varchar(255)"`
	Location string    `json:"location"`
	OrganizerID uint    `json:"organizer_id"`
	ImageURL string    `json:"image_url"`
	StartDate time.Time `json:"start_date"`

	Group Group `gorm:"foreignKey:GroupID"`
	Event Event `gorm:"foreignKey:EventID"`
}

func (gE *GroupEvent) BeforeCreate(tx *gorm.DB) error {
	gE.ID = uuid.NewString()
	return nil
}
