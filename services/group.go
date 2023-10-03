package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/types"
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

func UnsubscribeUserFromGroup(userID string, groupID string) error {
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

func ListGroups(limit int, offset int) ([]types.GroupDetailResponse, *int64, error) {
	var groupDetailsList []types.GroupDetailResponse
	var totalGroups int64

	query1 := `
       SELECT
    g.id AS id,
    g.name AS name,
    g.image AS image,
    g.created_at AS created_at,
    g.updated_at AS updated_at,
    COUNT(DISTINCT ge.id) AS events_count,
    COUNT(DISTINCT ug.id) AS members_count `
	query2 := "FROM `groups` g "

	query3 := `LEFT JOIN group_events ge ON g.id = ge.group_id
LEFT JOIN user_groups ug ON g.id = ug.group_id
GROUP BY
    g.id,
    g.name,
    g.image,
    g.created_at,
    g.updated_at`

	dbQuery := db.DB.Raw(query1 + query2 + query3)
	err := dbQuery.Model(&models.Group{}).Offset(offset).Limit(limit).Scan(&groupDetailsList).Error

	if err != nil {
		return nil, nil, err
	}

	dbQuery.Model(&models.Group{}).Count(&totalGroups)
	return groupDetailsList, &totalGroups, nil
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

	result := db.DB.Where("id = ?", id).Preload("Events").Preload("Members.User").First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	return &group, nil
}

func GetGroupEvents(id string) (*models.Group, error) {
	var group models.Group
	fmt.Printf("group id %s", id)

	result := db.DB.Where("id = ?", id).Preload("Events", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(3)
		})
	}).Preload("Members.User").First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	return &group, nil
}

func GetGroupWithDetails(id string) (*types.GroupDetailResponse, int, error) {
	var groupDetails types.GroupDetailResponse
	var group models.Group

	err := db.DB.Model(&models.Group{}).Where("id = ?", id).First(&group).Error

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	var eventCount int64
	var memberCount int64

	db.DB.Model(&models.GroupEvent{}).Where("group_id = ?", id).Count(&eventCount)
	db.DB.Model(&models.UserGroup{}).Where("group_id = ?", id).Count(&memberCount)

	groupDetails.Name = group.Name
	groupDetails.Image = group.Image
	groupDetails.ID = group.ID
	groupDetails.CreatedAt = group.CreatedAt
	groupDetails.UpdatedAt = group.UpdatedAt
	groupDetails.EventsCount = eventCount
	groupDetails.MembersCount = memberCount

	return &groupDetails, http.StatusOK, nil

}
