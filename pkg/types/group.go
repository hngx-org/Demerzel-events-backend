package types

import (
	"demerzel-events/internal/models"
	"time"
)

type GroupDetailResponse struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Image        string       `json:"image"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	EventsCount  int64        `json:"events_count"`
	MembersCount int64        `json:"members_count"`
	Tags         []models.Tag `json:"tags"`
	CreatorId    string       `json:"creator_id"`
}
