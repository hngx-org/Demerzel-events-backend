package handlers

import (
	"time"

	"demerzel-events/internal/db"

	"github.com/gin-gonic/gin"
)

func SubscribeUserToGroupHandler(c *gin.Context) {
	type UserGroup struct {
		Id        string `json:"id"`
		UserID    string `json:"user_id"`
		GroupID   string `json:"group_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	groupID := c.Param("group_id")
	userID := c.Param("user_id")

	userGroupInstance := UserGroup{
		Id: "",
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

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User successfully subscribed to group",
		"data": userGroupInstance,
	})
}