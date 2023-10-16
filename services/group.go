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

func CreateGroup(requestBody *models.NewGroupReqBody, creatorId string) (*models.Group, error) {
	var newGroup models.Group
	newGroup.Name = requestBody.Name
	newGroup.Image = requestBody.Image
	newGroup.CreatorId = creatorId

	// preload the tags based on their IDs
	tags := []models.Tag{}
	if err := db.DB.Where(requestBody.Tags).Find(&tags).Error; err != nil {
		return nil, err
	}

	// Ensure at least one valid tag was added
	if len(tags) == 0 {
		return nil, errors.New("at least one valid tag must be added")
	}

	// add those tags to the new group to be created
	newGroup.Tags = tags

	if err := db.DB.Create(&newGroup).Error; err != nil {
		return nil, err
	}
	return &newGroup, nil
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

func ListGroups(name string, limit int, offset int) (*[]map[string]interface{}, *int64, error) {
	var groupDetailsList []models.Group

	var totalGroups int64

	query := db.DB.Model(&models.Group{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.
		Preload("Tags").
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Members").
		Offset(offset).
		Limit(limit).
		Find(&groupDetailsList).Error
	if err != nil {
		return nil, nil, err
	}

	var response []map[string]interface{}
	for _, group := range groupDetailsList {
		groupMap := map[string]interface{}{
			"id":            group.ID,
			"name":          group.Name,
			"image":         group.Image,
			"created_at":    group.CreatedAt,
			"updated_at":    group.UpdatedAt,
			"tags":          group.Tags,
			"members_count": len(group.Members),
			"events_count":  len(group.Events),
			"creator_id":    group.CreatorId,
		}
		response = append(response, groupMap)
	}

	query.Count(&totalGroups)
	return &response, &totalGroups, nil
}

func GetGroupsByUserId(userId string, limit int, offset int) ([]models.Group, *int64, error) {
	var groups []models.Group
	var totalGroups int64

	dbQuery := db.DB.
		Preload("Tags").
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

	err := db.DB.Model(&models.Group{}).Where("id = ?", id).Preload("Tags").First(&group).Error

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
	groupDetails.Tags = group.Tags
	groupDetails.CreatorId = group.CreatorId

	return &groupDetails, http.StatusOK, nil

}
