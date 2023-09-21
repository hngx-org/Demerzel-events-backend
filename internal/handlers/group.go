package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateGroup(ctx *gin.Context) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Invalid request body format: %s", err.Error()))
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name

	services.CreateGroup(&newGroup)

	response.Success(ctx, http.StatusCreated, "Group created successfully", map[string]any{"group": newGroup})
}

func SubscribeUserToGroup(c *gin.Context) {
	groupID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	data, err := services.SubscribeUserToGroup(user.Id, groupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "Invalid user or group ID. Please check the values and try again")
			return
		} else if err.Error() == "user already exists in group" {
			response.Error(c, http.StatusConflict, "User already subscribed to group")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to subscribe user to group")
		return
	}

	response.Success(c, http.StatusOK, "User successfully subscribed to group", data)
}

func UnsubscribeFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	rawUser, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusConflict, "error: unable to retrieve user from context")
		return
	}
	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusConflict, "error: invalid user type in context")
		return
	}

	err := services.DeleteUserGroup(user.Id, string(groupID))

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User is not subscribed to this group, no need to unsubscribe
			response.Error(c, http.StatusNotFound, "User not subscribed to this group")
			return
		}

		response.Error(c, http.StatusConflict, "Failed to unsubscribe user from group")
		return
	}

	response.Success(c, http.StatusOK, "User successfully unsubscribed to group", map[string]any{})

}
func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body format: %s", err.Error()))
		return
	}

	code, data, err := services.UpdateGroupService(db.DB, req, id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Group updated successfully", data)
}

// GetUserGroups returns all group this user belongs to
func GetUserGroups(c *gin.Context) {

	userGroup, err := services.GetGroupsByUserId(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Fetched all groups", userGroup)
	return
}
