package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
func CreateGroup(ctx *gin.Context) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		response.Error(ctx, err.Error())
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name

	_, err := services.CreateGroup(&newGroup)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, "Group Created Successfully", newGroup)
}

func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	_, data, err := services.UpdateGroupService(db.DB, req, id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "Group Name updated successfully", data)
}
