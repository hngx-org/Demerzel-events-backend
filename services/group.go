package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
)

// UnsubscribeUserFromGroup removes a user from a group by user and group IDs.
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
