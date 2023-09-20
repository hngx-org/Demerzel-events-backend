package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"

	"gorm.io/gorm"
)



func GetUserFromDB(id string) (*models.UserData, error) {
	
	// get user from db

	var user models.UserData
	result := db.DB.Where("id = ?", id).First(&user)
	
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil // Return nil when the user is not found
        }
        return nil ,result.Error // Return the actual error for other errors

    }

    return &user, nil
}



func UpdateUserByID(user *models.UserData) error {
    result := db.DB.Save(user) // Save updates to the user

    if result.Error != nil {
        return result.Error // Return the error if saving fails
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




