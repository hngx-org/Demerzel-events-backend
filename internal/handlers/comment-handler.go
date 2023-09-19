package handlers

import (
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
