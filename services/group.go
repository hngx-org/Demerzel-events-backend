package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func CreateGroup(group *models.Group) (*models.Group, error) {
	if err := db.DB.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func UpdateGroupService(tx *gorm.DB, req models.UpdateGroupRequest, id string) (int, models.Group, error) {
	group := models.Group{
		ID: id,
	}

	err := group.GetGroupByID(tx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusNotFound, group, fmt.Errorf("group with the specified id does not exist")
		}
		return http.StatusBadRequest, group, err
	}

	// check if a `name` was passed in the request body
	if req.Name != "" {
		group.Name = req.Name
	}

	// update group id
	err = group.UpdateGroupByID(tx)
	if err != nil {
		return http.StatusInternalServerError, group, err
	}

	return http.StatusOK, group, nil
}

func DeleteGroup(tx *gorm.DB, id string) error {

	db := tx.Delete(&models.Group{}, "group_id")
	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected < 1 {
		return fmt.Errorf("group with id=%s doesn't exist", id)
	}

	tx.Where("group_id = ?", id).Delete(&models.UserGroup{})

	return nil
}
