package routes

import (
	"demerzel-events/internal/handlers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	// Add routes and handler you are working on here
	// handler from Demerzel-events-backend/internal/handlers/comment-handler.go file
	router.POST("/events/:event_id/comments", handlers.CreateComment)
}