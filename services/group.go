package services

import (
	"fmt"
	"strings"

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
	var err error
	groups := make([]models.Group, 0)

	if len(strings.TrimSpace(name)) != 0 {
		result := db.DB.Where("name LIKE ?", fmt.Sprintf("%s%s%s", "%", name, "%")).Find(&groups)
		err = result.Error
	} else {
		result := db.DB.Find(&groups)
		err = result.Error
	}

	if err != nil {
		return nil, err
	}

	return groups, nil
}
