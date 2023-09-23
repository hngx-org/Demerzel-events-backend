package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GroupEventById(c *gin.Context) {

	id := c.Param("id")

	group := models.Group{
		ID: id,
	}

	result, err := group.GetGroupEvent(db.DB)

	if err != nil {
		response.Error(c, 500, "Can't process your request")
		return
	}

	res := make(map[string]interface{})
	res["Result"] = result

	response.Success(c, 200, "List of events", res)
}

func CreateEventHandler(c *gin.Context) {
	var input models.NewEvent

	// Error if JSON request is invalid
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Unable to parse payload: %s", err.Error()))
		return
	}

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

	input.CreatorId = user.Id

	// Check if description field is empty or is a string
	if input.Description == "" {
		response.Error(c, http.StatusBadRequest, "Desciption field is empty")
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

	if reflect.ValueOf(input.StartDate).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "StartDate is not a string")
		return
	}
	// Check if end_date field is empty or is a string
	if input.EndDate == "" {
		response.Error(c, http.StatusBadRequest, "EndDate field is empty")
		return
	}

	if reflect.ValueOf(input.EndDate).Kind() != reflect.String {
		response.Error(c, http.StatusBadRequest, "EndDate is not a string")
		return
	}

	createdEvent, err := models.CreateEvent(db.DB, &input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusCreated, "Event Created", map[string]interface{}{"event": createdEvent})

}

func GetEventHandler(c *gin.Context) {

	eventID := c.Param("eventid")

	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	event, err := models.GetEventByID(db.DB, eventID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusOK, "Event details fetched", map[string]*models.Event{"event": event})
}

// ListEventsHandler lists all events
func ListEventsHandler(c *gin.Context) {

	events, err := models.ListEvents(db.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusOK, "All Events", map[string]interface{}{
		"events": events,
	})
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

	userGroups, _, err := services.GetGroupsByUserId(user.Id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to get groups which user belongs to:"+err.Error())
		return
	}

	if len(userGroups) == 0 {
		events := make([]models.Event, 0)
		response.Success(c, http.StatusOK, "Friend Events", map[string]interface{}{
			"events": events,
		})

		return
	}

	var userGroupIds []string
	for _, group := range userGroups {
		userGroupIds = append(userGroupIds, group.ID)
	}

	events, err := models.ListEventsInGroups(db.DB, userGroupIds)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Unable to get events: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Events", map[string]interface{}{
		"events": events,
	})

	return
}
