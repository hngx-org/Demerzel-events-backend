package handlers

import (
	"bytes"
	"demerzel-events/dependencies/cloudinary"
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

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

func UploadFile(c *gin.Context) {
	uploadedFile, _ := c.FormFile("file")
	if uploadedFile == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files specified"})
		return
	}
	log.Println(uploadedFile.Filename)

	file, err := uploadedFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read upload:" + err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to upload file:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
