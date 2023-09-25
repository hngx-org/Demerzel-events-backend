package services

import (
	"errors"
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

func SubscribeUserToGroup(userID, groupID string) (*models.UserGroup, error) {
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

// query filter struct
type Filter struct {
	Search struct {
		Name string
	}
}

// get groups
func ListGroups(f Filter) ([]models.Group, error) {
	var err error
	groups := make([]models.Group, 0)

	args := []any{"%", f.Search.Name, "%"}

	if f.Search.Name != "" {
		result := db.DB.Where("name LIKE ?", fmt.Sprintf("%s%s%s", args...)).Preload("Events").Find(&groups)
		err = result.Error
	}

	if f.Search.Name == "" {
		result := db.DB.Preload("Events").Find(&groups)
		err = result.Error
	}

	if err != nil {
		return make([]models.Group, 0), err
	}

	return groups, nil
}

func GetGroupsByUserId(userId string) ([]models.Group, int, error) {
	if _, err := GetUserById(userId); err != nil {
		return nil, http.StatusNotFound, err
	}
	var groups []models.Group
	res := db.DB.
		Joins("JOIN user_groups ON groups.id = user_groups.group_id").
		Where("user_groups.user_id = ?", userId).
		Preload("Events").Find(&groups)

	if res.Error != nil {
		return nil, http.StatusNotFound, res.Error
	}

	return groups, http.StatusOK, nil
}
func DeleteGroup(tx *gorm.DB, id string) error {

	// Delete group with specified id.
	db := tx.Delete(&models.Group{}, "group_id")
	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected < 1 {
		return fmt.Errorf("group with id=%s doesn't exist", id)
	}

	return nil
}

func GetGroupById(id string) (*models.Group, error) {
	var group models.Group
	fmt.Printf("group id %s", id)

	result := db.DB.Where("id = ?", id).Preload("Members.User").Preload("Events").First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Group doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	return &group, nil
}
