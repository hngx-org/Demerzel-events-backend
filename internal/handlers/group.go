package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateGroup(c *gin.Context) {
	req := models.UpdateGroupRequest{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "failed to parse request", "error": err.Error()})
		return
	}

	code, data, err := services.UpdateGroupService(db.DB, req, id)
	if err != nil {
		c.JSON(code, gin.H{"status": code, "message": "failed to parse request", "error": err.Error()})
		return
	}

	c.JSON(code, gin.H{
		"status":  "success",
		"message": "Group Name updated successfully",
		"data":    data,
	})
}
