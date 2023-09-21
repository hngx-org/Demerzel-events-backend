package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add the handler function you are working on here
// to be used inside the routes in Demerzel-events-backend/internal/routes/comment-routes.go

func CreateComment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "request success",
	})
}

func UpdateComment(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusBadRequest, "An error occurred while creating account")
		return
	}

	user, ok := rawUser.(models.User)
	if !ok {
		response.Error(c, http.StatusBadRequest, "An error occurred while creating account")
		return
	}

	var updateReq models.UpdateComment
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	data, err := services.UpdateCommentById(updateReq, user.Id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment updated successfully", map[string]any{"comment": data})
}
