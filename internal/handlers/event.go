package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"math"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GetGroupEventsHandler(c *gin.Context) {

	id := c.Param("id")

	events, err := services.GetGroupEvents(id)

	if err != nil {
		response.Error(c, 500, "Can't process your request")
		return
	}

	response.Success(c, 200, "List of events", map[string]interface{}{"group": events})
}

func CreateEventHandler(c *gin.Context) {
	var input models.NewEvent

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	// Error if JSON request is invalid
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Unable to parse payload: %s", err.Error()))
		return
	}

	input.CreatorId = user.Id

	// Check if description field is empty or is a string
	if input.Description == "" {
		response.Error(c, http.StatusBadRequest, "Description field is empty")
		return
	}

	if reflect.ValueOf(input.Description).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Description is not a string")
		return
	}

	// Check if thumbnail field is empty or is a string
	if input.Thumbnail == "" {
		response.Error(c, http.StatusBadRequest, "Thumbnail field is empty")
		return
	}

	if reflect.ValueOf(input.Thumbnail).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Thumbnail is not a string")
		return
	}

	// Check if location field is empty or is a string
	if input.Location == "" {
		response.Error(c, http.StatusBadRequest, "Location field is empty")
		return
	}

	if reflect.ValueOf(input.Location).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Location is not a string")
		return
	}

	// Check if title field is empty or is a string
	if input.Title == "" {
		response.Error(c, http.StatusBadRequest, "Title field is empty")
		return
	}

	if reflect.ValueOf(input.Title).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "Title is not a string")
		return
	}

	// Check if start_time field is empty or is a string
	if input.StartTime == "" {
		response.Error(c, http.StatusBadRequest, "StartTime field is empty")
		return
	}

	if reflect.ValueOf(input.StartTime).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "StartTime is not a string")
		return
	}

	// Check if end_time field is empty or is a string
	if input.EndTime == "" {
		response.Error(c, http.StatusBadRequest, "EndTime field is empty")
		return
	}

	if reflect.ValueOf(input.EndTime).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "EndTime is not a string")
		return
	}

	// Check if start_date field is empty or is a string
	if input.StartDate == "" {
		response.Error(c, http.StatusBadRequest, "StartDate field is empty")
		return
	}

	if !helpers.IsValidDate(input.StartDate) {
		response.Error(c, http.StatusBadRequest, "Invalid StartDate. Should follow format 2023-09-21")
		return
	}

	if reflect.ValueOf(input.StartDate).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "StartDate is not a string")
		return
	}
	// Check if end_date field is empty or is a string
	if input.EndDate == "" {
		response.Error(c, http.StatusBadRequest, "EndDate field is empty")
		return
	}

	if !helpers.IsValidDate(input.EndDate) {
		response.Error(c, http.StatusBadRequest, "Invalid EndDate. Should follow format 2023-09-21")
		return
	}

	if reflect.ValueOf(input.EndDate).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "EndDate is not a string")
		return
	}

	createdEvent, code, err := services.CreateEvent(&input)

	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	sb, _, err := services.SubscribeUserToEvent(user.Id, createdEvent.Id)

	fmt.Println(sb)

	if err != nil {
		response.Success(c, http.StatusCreated, "Event Created but fail to add user", map[string]interface{}{"event": createdEvent})
		return
	}

	event, code, err := services.GetEventByID(createdEvent.Id)

	if err != nil {
		response.Success(c, http.StatusCreated, "Event Created", map[string]interface{}{"event": createdEvent})
		return
	}

	response.Success(c, code, "Event Created", map[string]interface{}{"event": event})

}

func GetEventHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	if eventID == "" {
		response.Error(c, http.StatusBadRequest, "Event ID is required")
		return
	}

	event, code, err := services.GetEventByID(eventID)

	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Event details fetched", map[string]interface{}{"event": event})
}

// ListEventsHandler lists all events
func ListEventsHandler(c *gin.Context) {
	startDate := c.Query("start_date")

	// Extract query parameters for pagination
	limit, offset, err := helpers.GetLimitAndOffset(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	events, totalEvents, code, err := services.ListEvents(startDate, *limit, *offset)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Events retrieved successfully", map[string]interface{}{
		"events":       events,
		"total_events": *totalEvents,
	})
}

func ListUpcomingEventsHandler(c *gin.Context) {
	events, err := services.ListUpcomingEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusOK, "Events retrieved successfully", map[string]interface{}{"events": events})
}

func DeleteEventHandler(c *gin.Context) {
	rawUser, exists := c.Get("user")
	eventId := c.Param("event_id")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	code, err := services.DeleteEvent(eventId, user.Id)

	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Event deleted successfully", nil)
}

func ListFriendsEventsHandler(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	userGroups, _, err := services.GetGroupsByUserId(user.Id, math.MaxInt64, 1) // A hack for now, not that I woulxh home and start fighting

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to get groups which user belongs to:"+err.Error())
		return
	}

	if len(userGroups) == 0 {
		events := make([]models.Event, 0)
		response.Success(c, http.StatusOK, "Friend Events", map[string]interface{}{"events": events})

		return
	}

	var userGroupIds []string
	for _, group := range userGroups {
		userGroupIds = append(userGroupIds, group.ID)
	}

	events, code, err := services.ListEventsInGroups(userGroupIds)

	if err != nil {
		response.Error(c, code, "Unable to get events: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Events", map[string]interface{}{"events": events})
}

func SubscribeUserToEvent(c *gin.Context) {
	eventID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	event, code, err := services.GetEventByID(eventID)
	if event == nil {
		response.Error(c, code, err.Error())
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, code, err = services.SubscribeUserToEvent(user.Id, eventID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	services.NotifyEventCreatorOnUserSubscription(event.CreatorId, event.Title, user.Name)

	response.Success(c, http.StatusOK, "User successfully subscribed to event", nil)
}

func UnsubscribeFromEvent(c *gin.Context) {
	eventID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}
	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	event, code, err := services.GetEventByID(eventID)
	if event == nil {
		response.Error(c, code, "Event does not exist")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	code, err = services.UnsubscribeUserFromEvent(user.Id, eventID)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	services.NotifyEventCreatorOnUserUnSubscription(event.CreatorId, event.Title, user.Name)

	response.Success(c, http.StatusOK, "User successfully unsubscribed to event", nil)
}

func GetUserEventSubscriptions(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}
	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	events, err := services.GetUserEventSubscriptions(user.Id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User event subscriptions retrieved", map[string]interface{}{"events": events})
}
