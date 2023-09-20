package handlers

import (
	"time"

	"demerzel-events/internal/db"
	"demerzel-events/internal/models"

	"github.com/gin-gonic/gin"
)

func SubscribeUserToGroupHandler(c *gin.Context) {
	groupID := c.Param("group_id")
	userID := c.Param("user_id")

	userGroupInstance := models.UserGroup{
		UserID: string(userID),
		GroupID: string(groupID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := db.DB.Create(&userGroupInstance)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to subscribe user to group",
			"data": result.Error,
		})
		return
	}

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "User successfully subscribed to group",
		"data": userGroupInstance,
	})
}