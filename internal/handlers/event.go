package handlers

import (
	"bytes"
	"demerzel-events/dependencies/cloudinary"
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GroupEventById(c *gin.Context) {
	
	id := c.Param("id")

    group := models.Group{
        Id: id,	
	}

	result, err := group.GetGroupEvent(db.DB)

	if err != nil {
		response.Error(c, 500, "Can't process your request")
		return
	}

	res := make(map[string]interface{}) 
	res["Message"] = result

	response.Success(c,200, "List of events",res)
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

	createdEvent, err := models.CreateEvent(db.DB, &input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusCreated, "Event Created", map[string]interface{}{"event": createdEvent})

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

func UploadFileHandler(c *gin.Context) {
	uploadedFile, _ := c.FormFile("file")
	if uploadedFile == nil {
		response.Error(c, http.StatusBadRequest, "No files specified")
		return
	}

	file, err := uploadedFile.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "Unable to read upload:"+err.Error())
		return
	}

	uploader := cloudinary.Config{
		ApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		ApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		BaseUrl:   os.Getenv("CLOUDINARY_BASE_URL"),
	}
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), uploadedFile.Filename)

	url, err := uploader.UploadFile(buf.Bytes(), filename)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Unable to upload file:"+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "File uploaded", map[string]string{"url": url})
	return

}
