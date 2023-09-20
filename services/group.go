package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
)

// Group Service interface
type Group interface {
	List(string) ([]models.Group, error)
}

// Represents a Group Service
type group struct{}

// Creates a new Group service
func NewGroupService() Group {
	return &group{}
}

// Lists all groups
func (s *group) List(name string) ([]models.Group, error) {
	groups := make([]models.Group, 0)
	query := "SELECT * FROM groups WHERE lower(name) = lower(?) OR ? = '';"

	result := db.DB.Raw(query, name, name).Scan(&groups)
	if err := result.Error; err != nil {
		return []models.Group{}, err
	}

	return groups, nil
}
