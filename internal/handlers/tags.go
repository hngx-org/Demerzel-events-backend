package handlers

import (
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	tags, error := services.GetTags()
	if error != nil {
		response.Error(c, http.StatusInternalServerError, "Error getting list of tags")
		return
	}

	response.Success(c, http.StatusOK, "Tags successfully retrieved", map[string]interface{}{
		"tags": tags,
	})
}
