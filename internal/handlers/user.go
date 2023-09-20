package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateUserProfile(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
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
	// Get the user by ID
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error: cannot find user",
			"status": "error",
		})
		return
	}

	// Update the user's information
	user.Name = RequestBody.Name
	user.Email = RequestBody.Avatar

	// Save the updated user to the database
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error: could not save user to database",
			"status": "error",
		})
		return
	}
}