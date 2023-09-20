package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    ".../internal/models"
)

func CreateEvent(c *gin.Context) {
    // Check if the user is authenticated and authorized to create events
    if !isUserAuthenticatedAndAuthorized(c) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var event models.Event
    if err := c.ShouldBindJSON(&event); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set the organizer ID based on the authenticated user
    user := getCurrentUser(c)
    event.OrganizerID = user.ID

    if err := db.Create(&event).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, event)
}

func ListEvents(c *gin.Context) {
    var events []models.Event
    if err := db.Find(&events).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, events)
}

func GetEventByID(c *gin.Context) {
    eventID := c.Param("id") // If ID is the URL parameter for the event ID
    var event models.Event
    if err := db.Where("id = ?", eventID).First(&event).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
        return
    }
    c.JSON(http.StatusOK, event)
}
