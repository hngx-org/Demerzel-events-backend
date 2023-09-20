package handlers

import (
	"demerzel-events/internal/db"
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
func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "failed to parse request", "error": err.Error()})
		return
	}

	code, data, err := services.UpdateGroupService(db.DB, req, id)
	if err != nil {
		c.JSON(code, gin.H{"status": code, "message": "failed to parse request", "error": err.Error()})
		return
	}

	c.JSON(code, gin.H{
		"status":  "success",
		"message": "Group Name updated successfully",
		"data":    data,
	})
}

// List all groups
func ListGroups(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	f := services.Filter{
		Search: struct{ Name string }{
			Name: name,
		},
	}

	groups, err := services.ListGroups(f)
	if err != nil {
		response.Error(c, "error: failed to fetch groups")
		return
	}

	response.Success(c, "Groups retrieved successfully", map[string]any{
		"groups": groups,
	})
}
