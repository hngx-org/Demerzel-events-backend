package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpUser(user models.User) (models.UserResponse, string, int, error) {
	// check if user already exists
	_, err := getUserFromDB(user.Email)
	if err == nil {
		return models.UserResponse{}, "user already exist", 403, errors.New("user already exist in database")
	}

	// logic to signup user

	return models.UserResponse{}, "successfully created user", 0, nil
}

func getUserFromDB(email string) (models.User, error) {
	
	// get user from db
	var user models.User
	err:=db.DB.Preload("Events").Where("Email=?",email).Find(&user).Error
	if err != nil {
        return models.User{}, err
    }
	return user, nil
}

func UpdateUser(user *models.User, id string, c *gin.Context, RequestBody models.UserUpdate) {
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