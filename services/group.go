package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
)

func CreateGroup(group *models.Group) (*models.Group, error) {
	if err := db.DB.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}
