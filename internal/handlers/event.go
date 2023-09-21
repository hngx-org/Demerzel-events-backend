package handlers

import (
	"demerzel-events/internal/db"

	"demerzel-events/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GroupEvent(c *gin.Context) {

	//id := c.Param("id")
	//group, err := models.Group.Get(id)
	// if err != nil{

	// 	c.JSON(404, gin.H{
	// 		"status":  "Failed",
	// 		"message": "Group Not Found",
	// 	})
	// }
	
    //result, err := group.GetGroupEvent(db.DB)

	//this is for testing
	id := c.Param("id")

	if id == "" {
			c.JSON(400, gin.H{
			"status":  "Bad request",
			"message": "Invalid group",
		})
		return
	}

    group := models.Group{
        Id: id,	
	}

	result:= group.GetGroupEvent(db.DB)

		c.JSON(200, gin.H{
			"status":  "success",
			"message": result,
		})
		
}

func CreateEventHandler(c *gin.Context) {
	var input models.NewEvent

	// Error if JSON request is invalid
	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO use auth middleware to get creator of event.
	if input.CreatorId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id could not be found"})
		return
	}

	createdEvent, err := models.CreateEvent(db.DB, &input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Event Created": createdEvent})

}

// list all events
func ListEventsHandler(c *gin.Context) {

	events, err := models.ListEvents(db.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"events": events,
		},
	})
}
