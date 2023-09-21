package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"

	"gorm.io/gorm"
)

func CreateNewComment(comment *models.Comment) (*models.Comment, error) {
	if err := db.DB.Create(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
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
