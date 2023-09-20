package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/pkg/types"
	"demerzel-events/services"
	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, "error: unable to retrieve user from context")
		return
	}

	user, ok := rawUser.(models.User)
	if !ok {
		response.Error(c, "error: invalid user type in context")
		return
	}

	var updateData types.UserUpdatables
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.Error(c, "error: invalid json body format")
		return
	}

	services.UpdateUserById(&user, updateData)
}

func GetUserById(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		response.Error(c, "error: user id cannot be empty")
		return
	}

	user, err := services.GetUserById(id)
	if err != nil {
		response.Error(c, "error: unable to retrieve user")
		return
	}

	response.Success(c, "User retrieved successfully", map[string]any{"user": user})
}
func GetUsers(c *gin.Context){

	users, err := services.GetUsers()
	if err != nil {
		response.Error(c, "error: unable to retrieve users")
		return
	}
	response.Success(c, "Users Retrieved Successfully",map[string]any{"user": users})
}
