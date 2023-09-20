package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"demerzel-events/services"
)

type Group struct {
	Service services.Group
}

func (h *Group) List(c *gin.Context) {
	name, ok := c.GetQuery("name")
	if !ok {
		name = ""
	}

	groups, err := h.Service.List(name)
	if err != nil {
		c.JSON(500, map[string]string{
			"detail": "internal server error",
		})
		return
	}

	var message string

	if len(groups) == 0 {
		message = "no groups"
	} else {
		message = "Group(s) successfully fetched"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data":    groups,
	})
}
