package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

func CreateNewComment(newComment *models.NewComment, user *models.User) (*models.CommentResponse, error) {
	comment := models.Comment{
		Body:      newComment.Body,
		Images:    newComment.Images,
		EventId:   newComment.EventId,
		CreatorId: user.Id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	commentRespnse := &models.CommentResponse{
		Id:        comment.Id,
		Body:      comment.Body,
		Images:    comment.Images,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		EventId:   comment.EventId,
		Creator: models.CommentCreator{
			Id:     user.Id,
			Name:   user.Name,
			Avatar: user.Avatar,
		},
	}
	return commentRespnse, nil
}


func UpdateCommentById(updateReq *models.UpdateComment, userId string) (*models.Comment, error) {
	var comment *models.Comment
	result := db.DB.Where("id = ?", updateReq.Id).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("coomment doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	if comment.CreatorId != userId {
		return comment, errors.New("you are not authorized to update this comment")
	}

	comment.Body = updateReq.Body
	if err := db.DB.Save(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}

func GetCommentByCommentId(commentId string) (*models.CommentResponse, error) {
	var comment *models.Comment
	err := db.DB.Where("id = ?", commentId).Preload("Creator").First(&comment).Error
	// err := db.DB.Where("id = ?", commentId).Where("event_id = ?", eventId).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment does not exist")
		}
		return nil, err
	}

	commentResponse := &models.CommentResponse{
		Id:        comment.Id,
		Body:      comment.Body,
		Images:    comment.Images,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		EventId:   comment.EventId,
		Creator: models.CommentCreator{
			Id:     comment.Creator.Id,
			Name:   comment.Creator.Name,
			Avatar: comment.Creator.Avatar,
		},
	}

	return commentResponse, nil
}

func GetComments(eventId string, perPage, offset int) ([]*models.CommentResponse, int, error) {
	var comments []*models.Comment
	var totalComments int64

	// Query for comments with pagination
	err := db.DB.Where("event_id = ?", eventId).Preload("Creator").
		Offset(offset).Limit(perPage).Find(&comments).Error
	if err != nil {
		log.Println("Error fetching comments from db")
		return nil, int(totalComments), err
	}

	// Get the total count of comments for the event
	db.DB.Model(&models.Comment{}).Where("event_id = ?", eventId).Count(&totalComments)

	// Create a slice to hold the CommentResponse objects
	commentResponses := make([]*models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponse := &models.CommentResponse{
			Id:        comment.Id,
			Body:      comment.Body,
			EventId:   comment.EventId,
			Images:    comment.Images,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}

		// Check if comment.Creator is not nil and populate CommentCreator
		if comment.Creator != nil {
			commentResponse.Creator = models.CommentCreator{
				Id:     comment.Creator.Id,
				Name:   comment.Creator.Name,
				Avatar: comment.Creator.Avatar,
			}
		}

		commentResponses[i] = commentResponse
	}

	return commentResponses, int(totalComments), nil
}

func UpdateCommentById(updateReq *models.UpdateComment, userId string) (*models.Comment, error) {
	var comment *models.Comment
	result := db.DB.Where("id = ?", updateReq.Id).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("coomment doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	if comment.CreatorId != userId {
		return comment, errors.New("you are not authorized to update this comment")
	}

	comment.Body = updateReq.Body
	if err := db.DB.Save(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}

func DeleteCommentById(commentId string, userId string) error {
	var comment models.Comment
	result := db.DB.Where("id = ?", commentId).First(&comment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("coomment doesn't exist")
		}
		return result.Error // Return the actual error for other errors
	}

	if comment.CreatorId != userId {
		return errors.New("you are not authorized to delete this comment")
	}

	if err := db.DB.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}
