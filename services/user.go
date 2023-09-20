package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"

	"gorm.io/gorm"
)



func GetUserFromDB(id string) (models.User, error) {
	
	// get user from db
	var user models.User
	err:=db.DB.Preload("Events").Where("id=?",id).Find(&user).Error
	if err != nil {
        return models.User{}, err
    }
	return user, nil
}



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

