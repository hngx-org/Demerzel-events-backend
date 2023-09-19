package handlers

import (
	"demerzel-events/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateGroup(ctx *gin.Context, db *gorm.DB) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var newGroup models.Group
	newGroup.CreatedAt = time.Now()
	newGroup.UpdatedAt = time.Now()
	newGroup.Name = requestBody.Name

	db.Create(&newGroup)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "", "data": newGroup})
}
