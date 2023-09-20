package handlers

import (
	"github.com/gin-gonic/gin"

	"demerzel-events/pkg/response"
	"demerzel-events/services"
)

// List all groups
func ListGroups(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	f := services.Filter{
		Search: struct{ Name string }{
			Name: name,
		},
	}

	groups, err := services.ListGroups(f)
	if err != nil {
		response.Error(c, "error: failed to fetch groups")
		return
	}

	response.Success(c, "Groups retrieved successfully", response.ResponseData{
		"groups": groups,
	})
}
