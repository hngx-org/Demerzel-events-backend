package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
)

// Represents a filter whose values are used to filter query results
type Filter struct {
	Search struct {
		Name string
	}
}

func ListGroups(f Filter) ([]models.Group, error) {
	var err error
	groups := make([]models.Group, 0)

	if f.Search.Name != "" {
		result := db.DB.Where("name LIKE ?", f.Search.Name).Find(&groups)
		err = result.Error
	}

	if f.Search.Name == "" {
		result := db.DB.Find(&groups)
		err = result.Error
	}

	if err != nil {
		return make([]models.Group, 0), err
	}

	return groups, nil
}
