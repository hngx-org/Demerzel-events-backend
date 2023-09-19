package models

type User struct {
	// Add user fields

	// Events Relationship
	Events           []Event `gorm:"foreignKey:Creator"`
	InterestedEvents []Event `gorm:"many2many:interested_events;"`
}
