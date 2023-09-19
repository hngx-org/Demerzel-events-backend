package handlers

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserGroups returns all user groups by id
func GetUserGroups(c *gin.Context) {

	var userGroup []models.UserGroup
	userId := c.Param("id")

	res := db.DB.Find(&userGroup, "user_id = ?", userId)

	if res.Error != nil || res.RowsAffected <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": fmt.Sprintf("invalid or malformed user id"),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Fetched all groups",
		"data":    userGroup,
	})
	return
}
