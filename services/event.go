package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"gorm.io/gorm"
)

func parseEventResponse(event *models.Event) models.EventResponse {
	// Doing this because we don't need to return the whole comments field here
	// There has to be a cleaner and efficient way to do this
	var evResponse models.EventResponse
	evResponse.Id = (*event).Id
	evResponse.CreatorId = (*event).CreatorId
	evResponse.Thumbnail = (*event).Thumbnail
	evResponse.Location = (*event).Location
	evResponse.Title = (*event).Title
	evResponse.Description = (*event).Description
	evResponse.StartDate = (*event).StartDate
	evResponse.EndDate = (*event).EndDate
	evResponse.StartTime = (*event).StartTime
	evResponse.EndTime = (*event).EndTime
	evResponse.CreatedAt = (*event).CreatedAt
	evResponse.UpdatedAt = (*event).UpdatedAt
	evResponse.Creator = (*event).Creator
	evResponse.Groups = (*event).Groups

	customComments := []models.CommentCreator{}
	for _, comment := range event.Comments {
		var customComment models.CommentCreator
		customComment.Avatar = comment.Creator.Avatar
		customComment.Id = comment.Creator.Id
		customComment.Name = comment.Creator.Name

		customComments = append(customComments, customComment)
	}

	evResponse.Comments = customComments
	return evResponse
}

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

func GetEventByID(eventID string) (*models.EventResponse, int, error) {
	var event models.Event

	err := db.DB.Model(&models.Event{}).
		Where("id = ?", eventID).
		Preload("Creator").
		Preload("Attendees").
		Preload("Groups").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Limit(3).Order("created_at desc").Preload("Creator")
		}).
		First(&event).Error

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	// Doing this because we don't need to return the whole comments field here
	// There has to be a cleaner and efficient way to do this
	evResponse := parseEventResponse(&event)

	return &evResponse, http.StatusOK, nil
}

func CreateEvent(event *models.NewEvent) (*models.EventResponse, int, error) {
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
func ListEvents(startDate string, title string, dayOfWeek int, month int, weekNumber int, limit int, offset int) (*[]models.EventResponse, *int64, int, error) {
	var events []models.Event
	var totalEvents int64

	query := db.DB
	if startDate != "" && helpers.IsValidDate(startDate) {
		query = query.Where(&models.Event{StartDate: startDate})
	}

	// Filter by search query (event title)
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// Filter by day of the week (1 - 6 = Monday - Sunday)
	if dayOfWeek >= 1 && dayOfWeek <= 7 {
		fmt.Println("\nHHH", dayOfWeek)
		query = query.Where("DAYOFWEEK(start_date) = ?", dayOfWeek)
	}

	// Filter by month (1 - 12 = January - December)
	if month >= 1 && month <= 12 {
		query = query.Where("MONTH(start_date) = ?", month)
	}

	// Filter by week number - (There are 52 or 53 weeks in a year)
	if weekNumber >= 1 && weekNumber <= 53 {
		query = query.Where("WEEK(start_date) = ?", weekNumber)
	}

	err := query.
		Order("start_date, start_time").
		Preload("Creator").
		Preload("Attendees").
		Preload("Groups").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Limit(3).Order("created_at desc").Preload("Creator")
		}).
		Offset(offset).Limit(limit).Find(&events).Error

	var eventsResponse []models.EventResponse
	for _, event := range events {
		parsedEvent := parseEventResponse(&event)
		eventsResponse = append(eventsResponse, parsedEvent)
	}

	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	query.Model(&models.Event{}).Count(&totalEvents)

	if eventsResponse == nil {
		eventsResponse = []models.EventResponse{}
	}

	return &eventsResponse, &totalEvents, http.StatusOK, nil
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

func GetUserEventSubscriptions(userID string) (*[]models.Event, error) {
	var user models.User
	result := db.DB.Where("id = ?", userID).Preload("InterestedEvents").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user.InterestedEvents, nil
}

func GetEventAttendees(eventId string) (*[]models.User, error) {
	var event models.Event

	err := db.DB.Model(&models.Event{}).
		Where("events.id = ?", eventId).
		Select("events.id").
		Preload("Attendees").
		First(&event).Error

	if err != nil {
		return nil, err
	}

	return &event.Attendees, nil
}

func UpdateEvent(eventId string, userId string, data models.UpdateEvent) (*models.Event, error) {
	var event models.Event
	err := db.DB.Where("id = ?", eventId).First(&event).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event does not exist")
		}
		return nil, err
	}

	if event.CreatorId != userId {
		return nil, errors.New("user cannot update event")
	}

	updatable := []string{
		"Thumbnail",
		"Location",
		"Description",
		"StartDate",
		"EndDate",
		"StartTime",
		"EndTime",
	}

	for _, value := range updatable {
		field := reflect.ValueOf(&data).Elem().FieldByName(value)
		if field.IsValid() {
			if !field.IsZero() {
				reflect.ValueOf(&event).Elem().FieldByName(value).SetString(field.String())
			}
		}
	}

	err = db.DB.Save(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}
