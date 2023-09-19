package models

import "time"

type Event struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Creator     string    `json:"creator"`
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
	UserId  string `json:"user_id"`
	EventId string `json:"event_id"`

	User  User  `gorm:"foreignKey:User"`
	Event Event `gorm:"foreignKey:Event"`
}

type GroupEvent struct {
	GroupId string `json:"group_id"`
	EventId string `json:"event_id"`

	// needs Group model.
	Event Event `gorm:"foreignKey:Event"`
}
