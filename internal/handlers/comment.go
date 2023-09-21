package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Add the handler function you are working on here
// to be used inside the routes in Demerzel-events-backend/internal/routes/comment-routes.go

func CreateComment(c *gin.Context) {
	eventId := c.Param("event_id")
	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, "error: unable to retrieve user from context")
		return
	}

	// user, ok := rawUser.(models.User)

	fmt.Printf("\n1111 rawUser type: %T, rawUser value: %v\n", rawUser, rawUser)

	// if !ok {
	// 	response.Error(c, "error: invalid user type in context")
	// 	return
	// }

	var newComment models.Comment

	if err := c.BindJSON(&newComment); err != nil {
		response.JSON(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if strings.TrimSpace(newComment.Body) == "" {
		response.JSON(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// newComment.UserId = user.Id
	newComment.UserId = "d5d5416d-1b59-46b2-9ee3-38a3bb4c022b"
	newComment.CreatedAt = time.Now()
	newComment.UpdatedAt = time.Now()
	newComment.EventId = eventId
	newComment.BeforeCreate(db.DB)

	fmt.Println("1111", newComment, newComment.Images)

	data, err := services.CreateNewComment(&newComment)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "Comment created", map[string]any{"comment": data})
}

func UpdateComment(c *gin.Context) {
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

	var updateReq models.UpdateComment
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.JSON(c, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	data, err := services.UpdateCommentById(updateReq, user.Id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, "Comment updated successfully", map[string]any{"comment": data})

}
