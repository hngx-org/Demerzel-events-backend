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

func parseGroupResponse(group *models.Group) models.GroupResponse {
	// Doing this because we don't need to return the whole comments field in the events here
	// There has to be a cleaner and efficient way to do this
	var groupResponse models.GroupResponse
	groupResponse.ID = (*group).ID
	groupResponse.Name = (*group).Name
	groupResponse.Image = (*group).Image
	groupResponse.CreatedAt = (*group).CreatedAt
	groupResponse.UpdatedAt = (*group).UpdatedAt
	groupResponse.Members = (*group).Members

	groupEvents := []models.EventResponse{}
	for _, event := range (*group).Events {
		singleEvent := parseEventResponse(&event)
		groupEvents = append(groupEvents, singleEvent)
	}

	groupResponse.Events = groupEvents
	return groupResponse
}

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

func ListGroups(name string, limit int, offset int) ([]types.GroupDetailResponse, *int64, error) {
	var groupDetailsList []types.GroupDetailResponse
	var totalGroups int64

	query := db.DB.Model(&models.Group{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query = query.Select(`
		groups.id AS id,
		groups.name AS name,
		groups.image AS image,
		groups.created_at AS created_at,
		groups.updated_at AS updated_at,
		COUNT(DISTINCT group_events.id) AS events_count,
		COUNT(DISTINCT user_groups.id) AS members_count
	`).
		Joins("LEFT JOIN group_events ON groups.id = group_events.group_id").
		Joins("LEFT JOIN user_groups ON groups.id = user_groups.group_id").
		Group("groups.id, groups.name, groups.image, groups.created_at, groups.updated_at")

	err := query.Offset(offset).Limit(limit).Scan(&groupDetailsList).Error

	if err != nil {
		return nil, nil, err
	}

	query.Count(&totalGroups)
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

func GetGroupEvents(id string) (*models.GroupResponse, error) {
	var group models.Group
	fmt.Printf("group id %s", id)

	result := db.DB.
		Where("id = ?", id).
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at desc").Preload(("Creator")).Limit(3)
			})
		}).
		Preload("Members.User").
		First(&group)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group doesn't exist")
		}
		return nil, result.Error // Return the actual error for other errors
	}

	groupResponse := parseGroupResponse(&group)
	return &groupResponse, nil
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
