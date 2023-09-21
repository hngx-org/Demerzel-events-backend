package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/types"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil when the user is not found
		}
		return nil, result.Error // Return the actual error for other errors
	}

	return &user, nil
}

func GetUserById(id string) (*models.User, error) {
	var user models.User
	fmt.Printf("user id %s", id)

	result := db.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil when the user is not found
		}
		return nil, result.Error // Return the actual error for other errors
	}

	return &user, nil
}

func CreateUser(user *models.User) (*models.User, error) {
	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUserById(user *models.User, data types.UserUpdatables) error {
	user.Name = data.Name
	if data.Avatar != "" {
		user.Avatar = data.Avatar
	}

	if err := db.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}
