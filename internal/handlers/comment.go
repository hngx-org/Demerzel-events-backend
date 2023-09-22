package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Add the handler function you are working on here
// to be used inside the routes in Demerzel-events-backend/internal/routes/comment-routes.go

// To-Do
// Check if event id is valid before creating comment
func CreateComment(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	var input models.NewComment

	if err := c.BindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if strings.TrimSpace(input.Body) == "" {
		response.Error(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	data, err := services.CreateNewComment(&input, user.Id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment created", map[string]any{"comment": data})
}

func UpdateComment(c *gin.Context) {
	commentId := c.Param("comment_id")
	var updateReq models.UpdateComment

	updateReq.Id = commentId

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	updateReq.Body = strings.TrimSpace(updateReq.Body)

	data, err := services.UpdateCommentById(&updateReq, user.Id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment updated successfully", map[string]any{"comment": data})
}

func GetComment(c *gin.Context) {
	// eventId := c.Param("eventid")
	commentId := c.Param("comment_id")

	_, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}
	comment, err := services.GetCommentByCommentId(commentId)

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment fetched successfully", map[string]*models.Comment{"comment": comment})
}

func DeleteComment(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	commentId := c.Param("comment_id")
	err := services.DeleteCommentById(commentId, user.Id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment deleted successfully", nil)
}
