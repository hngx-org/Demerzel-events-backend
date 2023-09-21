package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/pkg/types"
	"demerzel-events/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	var updateData types.UserUpdatables
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body format")
		return
	}

	services.UpdateUserById(&user, updateData)
	response.Success(c, http.StatusOK, "User updated successfully", nil)
}

func GetUserById(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "User ID cannot be empty")
		return
	}

	user, err := services.GetUserById(id)
	if user == nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("User with the ID (%s) does not exist", id))
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "An error occurred while retrieving user")
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", map[string]any{"user": user})
}
