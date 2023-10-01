package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"errors"
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

func SubscribeUserToGroup(userID, groupID string) (*models.UserGroup, error) {
	var user models.User
	var userGroup models.UserGroup

	result := db.DB.Where("group_id = ?", groupID).Where("user_id = ?", userID).First(&userGroup)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userGroup = models.UserGroup{
			UserID:  userID,
			GroupID: groupID,
		}

		result = db.DB.Create(&userGroup)
		if result.Error != nil {
			return nil, result.Error
		}

		// get user from database
		user.Id = userID
		if err := user.GetUserByID(db.DB); err != nil {
			return nil, err
		}
		// assign the user gotten to the UserGroup struct User field
		userGroup.User = user

		return &userGroup, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return nil, fmt.Errorf("user already subscribed to group")
}

func DeleteUserGroup(userID, groupID string) error {
	var userGroup models.UserGroup
	result := db.DB.Where("group_id = ?", groupID).Where("user_id = ?", userID).First(&userGroup)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not subscribed to group")
		}
		return result.Error
	}

	// Delete the UserGroup
	result = db.DB.Delete(&userGroup)
	return result.Error
}

func UpdateGroupById(id string, data *models.UpdateGroupRequest) (*models.Group, int, error) {

	group := &models.Group{}

	result := db.DB.Where("id = ?", id).First(&group)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, fmt.Errorf("group doesn't exist")
		}
		return nil, http.StatusInternalServerError, result.Error
	}

	group.Name = data.Name

	result = db.DB.Save(&group)

	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	return group, http.StatusOK, nil
}

// get groups
func ListGroups(nameFilter string, limit int, offset int) ([]models.Group, *int64, error) {
	groups := make([]models.Group, 0)
	var totalGroups int64

	args := []any{"%", nameFilter, "%"}

	dbQuery := db.DB
	if nameFilter != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%s%s%s", args...))
	}

	err := dbQuery.Preload("Events").
		Offset(offset).Limit(limit).Find(&groups).Error

	if err != nil {
		return nil, nil, err
	}

	dbQuery.Model(&models.Group{}).Count(&totalGroups)

	return groups, &totalGroups, nil
}

func GetGroupsByUserId(userId string, limit int, offset int) ([]models.Group, *int64, error) {
	var groups []models.Group
	var totalGroups int64

	dbQuery := db.DB.
		Joins("JOIN user_groups ON groups.id = user_groups.group_id").
		Where("user_groups.user_id = ?", userId)

	res := dbQuery.Preload("Events").
		Offset(offset).Limit(limit).Find(&groups)

	if res.Error != nil {
		return nil, nil, res.Error
	}

	dbQuery.Model(&models.Group{}).Count(&totalGroups)
	return groups, &totalGroups, nil
}

func DeleteGroup(id string) (int, error) {

	group, err := GetGroupById(id)

	if err != nil {
		return http.StatusBadRequest, err
	}

	err = db.DB.Delete(&group).Error

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil

}

func GetGroupById(id string) (*models.Group, error) {
	var group models.Group
	fmt.Printf("group id %s", id)

	result := db.DB.Where("id = ?", id).Preload("Members.User").Preload("Events").First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	return &group, nil
}
