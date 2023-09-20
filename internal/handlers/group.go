package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateGroup(ctx *gin.Context) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name

	services.CreateGroup(&newGroup)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "", "data": newGroup})
}

func UnsubscribeFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, "error: unable to retrieve user from context")
		return
	}
	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, "error: invalid user type in context")
		return
	}

	err := services.DeleteUserGroup(user.Id, string(groupID))

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
