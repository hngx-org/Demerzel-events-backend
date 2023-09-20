package services

import (
	"demerzel-events/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func CreateEventService(tx *gorm.DB, event *models.NewEvent) (*models.Event, error) {
	request := models.Event{
		Creator:     event.Creator,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}

	err := tx.Create(&request)

	fmt.Print(event)

	if err.Error != nil {
		fmt.Print(err)

		return &models.Event{}, err.Error
	}
	return &request, nil

}
