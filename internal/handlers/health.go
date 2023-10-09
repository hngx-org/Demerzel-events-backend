package handlers

import (
	"demerzel-events/pkg/response"
	"net/http"
	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	response.Success(c, http.StatusOK, "Team Demerzel Events API", nil)
}
