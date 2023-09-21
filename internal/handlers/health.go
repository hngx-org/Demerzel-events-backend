package handlers

import (
	"demerzel-events/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthHandler(c *gin.Context) {
	response.Success(c, http.StatusOK, "Team Demerzel Events API", nil)
}
