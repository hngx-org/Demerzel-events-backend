package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func AddReaction(user *models.User, eventId string, value string) error {
	event, _, err := GetEventByID(eventId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event does not exist")
		}
		return err
	}

	newReaction := models.NewReaction(user.Id, event.Id, value)
	err = db.DB.Create(&newReaction).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateReaction(user *models.User, reactionId string, value string) error {
	var reaction models.Reaction
	err := db.DB.Where("id = ?", reactionId).First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("reaction does not exist")
		}
		return err
	}

	if reaction.UserId != user.Id {
		return fmt.Errorf("user not allowed to update reaction")
	}

	reaction.Reaction = value
	err = db.DB.Save(&reaction).Error
	if err != nil {
		return err
	}

	return nil
}

func RemoveReaction(user *models.User, reactionId string) error {
	var reaction models.Reaction
	err := db.DB.Where("id = ?", reactionId).First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("reaction does not exist")
		}
		return err
	}

	if reaction.UserId != user.Id {
		return fmt.Errorf("user not allowed to remove reaction")
	}

	err = db.DB.Delete(&reaction).Error
	if err != nil {
		return err
	}

	return nil
}

func GetReactionsForEvent(eventId string, value string) (*[]models.Reaction, error) {
	var reactions []models.Reaction
	err := db.DB.
		Where("event_id = ?", eventId).
		Where("reaction = ?", value).
		Preload("User").
		Find(&reactions).Error

	if err != nil {
		return nil, err
	}

	return &reactions, nil
}

func GetAllReactionsForEvent(eventId string) (*[]models.Reaction, error) {
	var reactions []models.Reaction
	err := db.DB.
		Where("event_id = ?", eventId).
		Preload("User").
		Find(&reactions).Error

	if err != nil {
		return nil, err
	}

	return &reactions, nil
}

func GetReactionForEvent(userId string, eventId string, value *string) (*models.Reaction, error) {
	var reaction models.Reaction
	var err error

	if value != nil {
		err = db.DB.
			Where("event_id = ?", eventId).
			Where("user_id = ?", userId).
			Where("reaction = ?", value).
			First(&reaction).Error
	} else {
		err = db.DB.
			Where("event_id = ?", eventId).
			Where("user_id = ?", userId).
			First(&reaction).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &reaction, nil
}

func GetReactionById(reactionId string) (*models.Reaction, error) {
	var reaction models.Reaction
	err := db.DB.Where("id = ?", reactionId).First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &reaction, nil
}
