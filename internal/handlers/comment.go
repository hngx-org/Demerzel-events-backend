package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/helpers"
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

	_, code, eventErr := services.GetEventByID( input.EventId)

	if eventErr != nil {
		if eventErr.Error() == "record not found" {
			response.Error(c, code, "Event does not exist to comment on")
			return
		}

		response.Error(c, http.StatusInternalServerError, eventErr.Error())
		return
	}

	data, err := services.CreateNewComment(&input, user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment created", map[string]any{"comment": data})
}

func GetCommentsHandler(c *gin.Context) {
	// Extract query parameters for pagination
	limit, offset, err := helpers.GetLimitAndOffset(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	eventId := c.Param("event_id")

	_, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusBadRequest, "An error occurred")
		return
	}

	_, code, eventExist := services.GetEventByID( eventId)

	if eventExist != nil {
		response.Error(c, code, "Event does not exist")
		return
	}

	comments, totalComments, err := services.GetComments(eventId, *limit, *offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error Could not access the database")
		return
	}

	response.Success(c, http.StatusOK, "Comments Successfully retrieved", map[string]interface{}{
		"comments":       comments,
		"total_comments": totalComments,
	})
}

func UpdateComments(c *gin.Context) {
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
