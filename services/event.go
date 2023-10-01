package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func ListUpcomingEvents() ([]models.Event, error) {

	var events []models.Event

	now := time.Now().Format("2006-01-02 15:04:05")

	query := db.DB.Where("start_date > ?", now).Order("start_date, start_time")

	err := query.Preload("Creator").Preload("Attendees").Preload("Groups").Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil

}

func DeleteEvent(eventID string, userID string) (int, error) {

	event, code, err := GetEventByID(eventID)

	if err != nil {
		return code, err
	}

	if event.CreatorId != userID {
		return http.StatusUnauthorized, fmt.Errorf("you are not authorized to delete this event")
	}

	err = db.DB.Delete(&event).Error

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetEventByID(eventID string) (*models.Event, int, error) {
	var event models.Event

	err := db.DB.Where("id = ?", eventID).Preload("Creator").Preload("Attendees").Preload("Groups").First(&event).Error

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &event, http.StatusOK, nil
}

func CreateEvent(event *models.NewEvent) (*models.Event, int, error) {
	createdEvent := models.Event{
		CreatorId:   event.CreatorId,
		Thumbnail:   event.Thumbnail,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}

	result := db.DB.Create(&createdEvent)

	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(event.GroupId) >= 1 {
		for i := 0; i < len(event.GroupId); i++ {
			newGroupEvent := models.NewGroupEvent{
				EventId: createdEvent.Id,
				GroupId: event.GroupId[i],
			}
			_, code, err := CreateGroupEvent(&newGroupEvent)

			if err != nil {
				return nil, code, err
			}
		}
	}

	return GetEventByID(createdEvent.Id)
}

func CreateGroupEvent(groupEvent *models.NewGroupEvent) (*models.GroupEvent, int, error) {
	request := models.GroupEvent{
		EventId: groupEvent.EventId,
		GroupId: groupEvent.GroupId,
	}

	result := db.DB.Create(&request)

	if result.Error != nil {
		return &models.GroupEvent{}, http.StatusBadGateway, result.Error
	}

	return &request, http.StatusOK, nil
}

// ListEvents retrieves all events.
func ListEvents(startDate string, limit int, offset int) ([]models.Event, *int64, int, error) {
	var events []models.Event
	var totalEvents int64

	query := db.DB.Order("start_date, start_time")

	if startDate != "" && helpers.IsValidDate(startDate) {
		query.Where(&models.Event{StartDate: startDate})
	}
	err := query.Preload("Creator").Preload("Attendees").Preload("Groups").
		Offset(offset).Limit(limit).Find(&events).Error

	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	query.Model(&models.Event{}).Count(&totalEvents)

	return events, &totalEvents, http.StatusOK, nil
}

func ListEventsInGroups(groupIDs []string) ([]models.Event, int, error) {
	var events []models.Event

	err := db.DB.Model(&models.Event{}).
		Joins("JOIN group_events ON events.id = group_events.event_id").
		Where("group_events.group_id IN ?", groupIDs).
		Distinct("events.id").
		Select("events.id, events.creator_id, events.thumbnail," +
			"events.location, events.title, events.description, " +
			"events.start_date,events.end_date,events.start_time,events.end_time,events.created_at, events.updated_at").
		Preload("Creator").Find(&events).Error

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return events, http.StatusOK, nil
}

func SubscribeUserToEvent(userID, eventID string) (*models.InterestedEvent, int, error) {
	var interestedEvent models.InterestedEvent
	result := db.DB.Where("event_id = ?", eventID).Where("user_id = ?", userID).First(&interestedEvent)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		interestedEvent = models.InterestedEvent{
			UserId:  userID,
			EventId: eventID,
		}

		result = db.DB.Create(&interestedEvent)
		if result.Error != nil {
			return nil, http.StatusInternalServerError, result.Error
		}

		return &interestedEvent, http.StatusOK, nil
	}

	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	return nil, http.StatusForbidden, fmt.Errorf("user already subscribed to event")
}

func UnsubscribeUserFromEvent(userID, eventID string) (int, error) {
	var interestedEvent models.InterestedEvent
	result := db.DB.Where("event_id = ?", eventID).Where("user_id = ?", userID).First(&interestedEvent)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusForbidden, fmt.Errorf("user not subscribed to event")
		}
		return http.StatusInternalServerError, result.Error
	}

	// Delete the UserGroup
	result = db.DB.Delete(&interestedEvent)
	return http.StatusInternalServerError, result.Error
}

func GetUserEventSubscriptions(userID string) (*[]models.Event, int, error) {
	var user models.User
	result := db.DB.Where("id = ?", userID).Preload("InterestedEvents").First(&user)
	if result.Error != nil {
		return nil, http.StatusNotFound, result.Error
	}

	return &user.InterestedEvents, http.StatusOK, nil
}
