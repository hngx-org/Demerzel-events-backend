package handlers

import (
	"fmt"
	"net/http"

	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name

	services.CreateGroup(&newGroup)

	response.Success(
		ctx,
		http.StatusCreated,
		"Group created successfully",
		newGroup,
	)
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

	group, err := services.GetGroupById(groupID)
	if group == nil {
		response.Error(c, http.StatusNotFound, "Group does not exist")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := services.SubscribeUserToGroup(user.Id, groupID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
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

	group, err := services.GetGroupById(groupID)
	if group == nil {
		response.Error(c, http.StatusNotFound, "Group does not exist")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = services.DeleteUserGroup(user.Id, groupID)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User successfully unsubscribed to group", nil)
}

func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid request body format: %s", err.Error()),
		)
		return
	}

	code, data, err := services.UpdateGroupService(db.DB, req, id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}

	response.Success(c, code, "Group updated successfully", data)
}

func ListGroups(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	f := services.Filter{
		Search: struct{ Name string }{
			Name: name,
		},
	}

	groups, err := services.ListGroups(f)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error: failed to fetch groups")
		return
	}

	var message string

	if len(groups) == 0 {
		message = "No groups"
	}

	if len(groups) > 0 {
		message = "Groups retrieved successfully"
	}

	response.Success(c, http.StatusOK, message, groups)
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

	userGroups, code, err := services.GetGroupsByUserId(user.Id)
	if err != nil {
		response.Error(c, code, err.Error())
		return
	}
	response.Success(c, code, "Fetched all user groups", userGroups)
}

func GetGroupById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Group ID is required")
		return
	}

	group, err := services.GetGroupById(id)
	if group == nil {
		response.Error(c, http.StatusNotFound, "Group does not exist")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
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
	err := services.DeleteGroup(db.DB, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("group with id=%s deleted successfully", id), nil)

}
