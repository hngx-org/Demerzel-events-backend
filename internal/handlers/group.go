package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"

	"github.com/gin-gonic/gin"
)

func UnsubscribeUserFromGroupHandler(c *gin.Context) {
	groupID := c.Param("group_id")
	userID := c.Param("user_id")

	var userGroup models.UserGroup

	// result := db.DB.First(&userGroupInstance)
	// result := db.DB.Where("UserID = ? AND GroupID= ?", string(userID), string(groupID)).First(&userGroupInstance)
	result := db.DB.Where(&models.UserGroup{
		UserID:  string(userID),
		GroupID: string(groupID),
	}).First(&userGroup)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "User not subscribed to this group",
		})
		return
	}

	result = db.DB.Delete(&userGroup)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to unsubscribe user from group",
		})
		return
	}

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "User successfully unsubscribed to group",
	})
}
