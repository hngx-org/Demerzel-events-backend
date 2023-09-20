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

func UpdateCommentById(updateReq models.UpdateComment, userId string) (models.Comment, error) {
	var comment models.Comment
	result := db.DB.Where("event_id = ?", updateReq.EventId).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return comment, nil // Return nil when the user is not found
		}
		return comment, result.Error // Return the actual error for other errors
	}

	if comment.UserId != userId {
		return comment, errors.New("you are not authorized to update this comment")
	}

	comment.Body = updateReq.Body
	if err := db.DB.Save(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}
