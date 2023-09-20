package handlers

import (
	"demerzel-events/pkg/response"
	"demerzel-events/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UnsubscribeFromGroup(c *gin.Context) {
	groupID := c.Param("group_id")
	userID := c.Param("user_id")

	err := services.DeleteUserGroup(string(userID), string(groupID))

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User is not subscribed to this group, no need to unsubscribe
			response.Error(c, "User not subscribed to this group")
			return
		}

		response.Error(c, "Failed to unsubscribe user from group")
		return
	}

	response.Success(c, "User successfully unsubscribed to group", map[string]any{})
}
