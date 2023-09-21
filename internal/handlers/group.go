package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name

	services.CreateGroup(&newGroup)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "", "data": newGroup})
}

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

func DeleteGroupHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Please provide an id"})
		return
	}
	err := services.DeleteGroup(db.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Group with id not found", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": fmt.Sprintf("group with id=%s doesn't exist", id)})

}
