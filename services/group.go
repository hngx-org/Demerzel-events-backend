package services

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
)

func CreateGroup(group *models.Group) (*models.Group, error) {
	if err := db.DB.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func DeleteUserGroup(userID, groupID string) error {
	var userGroup models.UserGroup

	// Find the UserGroup by user and group IDs
	result := db.DB.Where(&models.UserGroup{
		UserID:  userID,
		GroupID: groupID,
	}).First(&userGroup)

	if result.Error != nil {
		return result.Error // Return the actual error for other errors
	}

	// Delete the UserGroup
	result = db.DB.Delete(&userGroup)
	return result.Error
}

func UpdateGroupService(
	tx *gorm.DB,
	req models.UpdateGroupRequest,
	id string,
) (int, models.Group, error) {
	group := models.Group{
		ID: id,
	}

	err := group.GetGroupByID(tx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusNotFound, group, fmt.Errorf(
				"group with the specified id does not exist",
			)
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

type Filter struct {
	Search struct {
		Name string
	}
}

func ListGroups(f Filter) ([]models.Group, error) {
	var err error
	groups := make([]models.Group, 0)

	if f.Search.Name != "" {
		result := db.DB.Where("name LIKE ?", f.Search.Name).Find(&groups)
		err = result.Error
	}

	if f.Search.Name == "" {
		result := db.DB.Find(&groups)
		err = result.Error
	}

	if err != nil {
		return make([]models.Group, 0), err
	}

	return groups, nil
}
