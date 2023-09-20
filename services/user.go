package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
 main

	"gorm.io/gorm"
)



func GetUserFromDB(id string) (models.User, error) {
	
	// get user from db
	var user models.User
	err:=db.DB.Preload("Events").Where("id=?",id).Find(&user).Error

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
 auth
	if err != nil {
        return models.User{}, err
    }
	return user, nil
}

 main


func UpdateUser(data models.UpdateUserStruct) error{
	user, err := GetUserFromDB(data.Id)
	if err != nil{
		return err
	}
	user.Email = data.Email
	user.Name = data.Name
	if err := db.DB.Save(user).Error; err != nil{
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*models.UserData, error) {
    var user models.UserData
    result := db.DB.Where("email = ?", email).First(&user)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil // Return nil when the user is not found
        }
        return nil, result.Error // Return the actual error for other errors
    }
    
    return &user, nil
}


func CreateUser(user *models.UserData) error {
    if err := db.DB.Create(user).Error; err != nil {
        return err
    }
    return nil
}
func UpdateUserByEmail(user *models.UserData) error {
    if err := db.DB.Where("email = ?", user.Email).Updates(user).Error; err != nil {
        return err
    }
    return nil
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
 auth
