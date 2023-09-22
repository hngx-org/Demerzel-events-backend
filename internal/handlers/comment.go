package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"
	"strconv"
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

	_, eventErr := models.GetEventByID(db.DB, input.EventId)

	if eventErr != nil {
		if eventErr.Error() == "record not found" {
			response.Error(c, http.StatusNotFound, "Event does not exist to comment on")
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
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "10")

	// Convert page and perPage parameters to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid per_page parameter")
		return
	}

	// Calculate the offset for pagination
	offset := (pageInt - 1) * perPageInt

	eventId := c.Param("event_id")

	_, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusBadRequest, "An error occurred while creating account")
		return
	}

	_, eventexist := models.GetEventByID(db.DB, eventId)
	if eventexist != nil {
		if eventexist.Error() == "record not found" {
			response.Error(c, http.StatusNotFound, "Event doesn't exist")
			return
		}
		response.Error(c, http.StatusInternalServerError, "An error occurred")
		return
	}

	comments, totalComments, err := services.GetComments(eventId, perPageInt, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error Could not access the database")
		return
	}

	response.Success(c, http.StatusOK, "Comments Successfully retrieved", map[string]interface{}{
		"comments":      comments,
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

func GetComment(c *gin.Context) {
	// eventId := c.Param("eventid")
	commentId := c.Param("comment_id")

	_, exists := c.Get("user")

	if !exists {
		response.Error(c, http.StatusInternalServerError, "An error occured")
		return
	}

	comment, commentErr := services.GetCommentByCommentId(commentId)

	if commentErr != nil {

		response.Error(c, http.StatusBadRequest, commentErr.Error())
		return
	}

	response.Success(c, http.StatusOK, "Comment fetched successfully", map[string]*models.CommentResponse{"comment": comment})
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
