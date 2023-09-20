package handlers

import (
	"demerzel-events/internal/db"

	"demerzel-events/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create event handler
func CreateEventHandler(c *gin.Context) {
	var input models.NewEvent

	// Error if JSON request is invalid
	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user id is null
	if input.Creator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id could not be found"})
		return
	}
	// I had an issue when creating events that it would create but it would not reference the EventCreator details, so what i did is to get the user by id, which is also a good error layer incase the user does not exist. I get the user information and append it to the Event details.
	// user, error := models.FindUserById(input.Creator, db.DB)
	// fmt.Print(error)
	// fmt.Print(user)
	// This logic has an issue, will refer back to it later
	createdEvent, error := models.CreateEvent(db.DB, &input)

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Event Created": createdEvent})

}
func ListEventsHandler(c *gin.Context) {
	// Your logic to fetch all events from the database or wherever they are stored.
	// This should include all events for everyone.

	// Example: Fetch all events and sort them by start date and start time using the models package.
	events, err := models.ListEvents(db.DB) // Assuming db.DB is your database connection.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
