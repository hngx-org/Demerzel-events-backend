package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	db := db.DB
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var newGroup models.Group
	newGroup.Name = requestBody.Name

	db.Create(&newGroup)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "", "data": newGroup})
}
