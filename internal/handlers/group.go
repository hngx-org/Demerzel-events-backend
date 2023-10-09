package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {

	rawUser, exists := ctx.Get("user")

	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(ctx, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	var requestBody struct {
		Name  string `json:"name" binding:"required"`
		Image string `json:"image" binding:"required"`
	}

	if err := ctx.ShouldBind(&requestBody); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name
	newGroup.Image = requestBody.Image

	group, err := services.CreateGroup(&newGroup)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	services.SubscribeUserToGroup(user.Id, group.ID)
	services.SendNewGroupNotificationToAllGroupNotificationEnabledUsers(newGroup.Name, user.Name, user.Id)

	response.Success(
		ctx,
		http.StatusCreated,
		"Group created successfully",
		newGroup,
	)
}

func SubscribeUserToGroup(c *gin.Context) {
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

	groupId := c.Param("id")
	group, err := services.GetGroupById(groupId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = services.SubscribeUserToGroup(user.Id, group.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User subscribed to group successfully", nil)
}

func UnsubscribeUserFromGroup(c *gin.Context) {
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

	groupId := c.Param("id")
	group, err := services.GetGroupById(groupId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = services.UnsubscribeUserFromGroup(user.Id, group.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User unsubscribed from group successfully", nil)
}

func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	data, code, err := services.UpdateGroupById(id, &req)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Group updated successfully", data)
}

func ListGroups(c *gin.Context) {
	name := c.Query("name")

	// Extract query parameters for pagination
	limit, offset, err := helpers.GetLimitAndOffset(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	groups, totalGroups, err := services.ListGroups(name, *limit, *offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error: failed to fetch groups")
		return
	}

	response.Success(c, http.StatusOK, "Groups retrieved successfully", map[string]interface{}{
		"groups":       groups,
		"total_groups": totalGroups,
	})
}

// GetUserGroups returns all group this user belongs to
func GetUserGroups(c *gin.Context) {
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

	// Extract query parameters for pagination
	limit, offset, err := helpers.GetLimitAndOffset(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userGroups, totalUserGroups, err := services.GetGroupsByUserId(user.Id, *limit, *offset)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Fetched all user groups", map[string]interface{}{
		"user_groups":       userGroups,
		"total_user_groups": totalUserGroups,
	})
}

func GetGroupById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Group ID is required")
		return
	}

	group, code, err := services.GetGroupWithDetails(id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Group retrieved successfully", group)
}

func DeleteGroup(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Please provide a valid group id")
		return
	}
	code, err := services.DeleteGroup(id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("group with id=%s deleted successfully", id), nil)

}

func GroupEventsById(c *gin.Context) {
	id := c.Param("id")

	result, code, err := services.GetEventByID(id)

	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Group events retrieved", result)
}
