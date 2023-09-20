package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUserProfile(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		c.JSON(500, gin.H{
			"message": "error: parameter parsed is not int",
			"status":  "error",
		})
	}

	var RequestBody models.UserUpdate
	if err := c.ShouldBindJSON(&RequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error: invalid json body format",
			"status": "error",
		})
		return
	}

	var user models.User

	services.UpdateUser(&user, id, c, RequestBody)
}