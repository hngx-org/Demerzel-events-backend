package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reaction struct {
	Id       string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserId   string `json:"creator_id" gorm:"type:varchar(255)"`
	EventId  string `json:"event_id" gorm:"type:varchar(255)"`
	Reaction string `json:"reaction" gorm:"type:varchar(20)"`

	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventId"`
}

func (u *Reaction) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.NewString()

	return nil
}

func NewReaction(userId string, eventId string, reaction string) *Reaction {
	return &Reaction{
		UserId:   userId,
		EventId:  eventId,
		Reaction: reaction,
	}
}
