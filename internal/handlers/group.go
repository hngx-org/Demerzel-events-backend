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
	}

	c.JSON(http.StatusOK, groups)
}
