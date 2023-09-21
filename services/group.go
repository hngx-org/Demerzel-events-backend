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

func SubscribeUserToGroup(userID, groupID string) (models.UserGroup, error) {
	var user models.User
	var group models.Group

	userGroup := models.UserGroup{
		UserID: userID,
		GroupID: groupID,
	}

	existingUser := db.DB.Where(&models.User{
		Id:  userID,
	}).First(&user)
	if existingUser.Error != nil {
		return userGroup, existingUser.Error
	}

	existingGroup := db.DB.Where(&models.Group{
		ID:  groupID,
	}).First(&group)
	if existingGroup.Error != nil {
		return userGroup, existingGroup.Error
	}

	userExistInGroup := db.DB.Where(&models.UserGroup{
		UserID:  userID,
	}).First(&userGroup)

	fmt.Println(userExistInGroup.Error)

	if userExistInGroup.Error == nil {
		return userGroup, fmt.Errorf("user already exists in group")
	}

	result := db.DB.Create(&userGroup)
	if result.Error != nil {
		return userGroup, result.Error
	}

	return userGroup, nil
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

func GetGroupsByUserId(userId string) ([]models.UserGroup, error) {
	if _, err := GetUserById(userId); err != nil {
		return nil, err
	}
	var userGroup []models.UserGroup
	res := db.DB.Find(&userGroup, "user_id = ?", userId)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return userGroup, nil
}
