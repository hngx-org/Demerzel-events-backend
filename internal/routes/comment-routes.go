package routes

import (
	"demerzel-events/internal/handlers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.RouterGroup) {
	// Add routes and handler you are working on here
	// handler from Demerzel-events-backend/internal/handlers/comment-handler.go file
	router.POST("/events/:event_id/comments", handlers.CreateComment)
	router.PUT("/events/:event_id/comments/:comment_id", handlers.UpdateComment)
}
