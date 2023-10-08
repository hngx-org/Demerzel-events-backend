package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/types"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func CreateNotification(notificationType, content string) (*models.Notification, int, error) {
	var notification models.Notification

	notification.Type = notificationType
	notification.Content = content

	err := db.DB.Create(&notification).Error

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &notification, http.StatusOK, nil
}

func GetNotificationByID(notificationID string) (*models.Notification, int, error) {
	var notification models.Notification

	err := db.DB.Where("id = ?", notificationID).First(&notification).Error

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &notification, http.StatusOK, nil
}

func ListNotifications() ([]models.Notification, int, error) {
	var notifications []models.Notification

	err := db.DB.Find(&notifications).Error

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return notifications, http.StatusOK, nil
}

func CreateUserNotification(userID, notificationID string) (*models.UserNotification, int, error) {
	var userNotification models.UserNotification

	userNotification.UserID = userID
	userNotification.NotificationID = notificationID

	err := db.DB.Create(&userNotification).Error

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &userNotification, http.StatusOK, nil
}

func GetUserNotificationByID(userNotificationID string) (*types.UserNotificationResponse, int, error) {
	var notificationsResponse types.UserNotificationResponse

	err := db.DB.Model(&models.UserNotification{}).Where("id = ?", userNotificationID).Preload("Notification").First(&notificationsResponse).Error

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &notificationsResponse, http.StatusOK, nil
}

func UpdateUserNotification(userNotificationID string, read bool) (int, error) {
	var userNotification models.UserNotification

	err := db.DB.Model(&models.UserNotification{}).Where("id = ?", userNotificationID).First(&userNotification).Error

	if err != nil {
		return http.StatusNotFound, err
	}

	userNotification.Read = read

	err = db.DB.Save(&userNotification).Error

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func ListUserNotifications(userID string) ([]types.UserNotificationResponse, int, error) {
	var notificationsResponse []types.UserNotificationResponse

	err := db.DB.Model(&models.UserNotification{}).Where("user_id = ?", userID).Preload("Notification").Find(&notificationsResponse).Error

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return notificationsResponse, http.StatusOK, nil

}

func SendNewGroupNotificationToAllGroupNotificationEnabledUsers(groupName, creatorName, creatorId string) (int, error) {

	enable := true

	users, _, _ := GetUsersWithEnabledNotificationSettings(&enable, nil, nil, nil)

	notification, _, err := CreateNotification("group", "New group "+groupName+" has been created"+" by "+creatorName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, user := range users {

		if user.Id == creatorId {
			continue
		}

		_, _, err = CreateUserNotification(user.Id, notification.ID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func SendNewEventNotificationToAllEventNotificationEnabledUsers(eventName, creatorName, creatorId string) (int, error) {
	
	enable := true

	users, _, _ := GetUsersWithEnabledNotificationSettings(nil, &enable, nil, nil)

	notification, _, err := CreateNotification("event", "New event "+eventName+" has been created"+" by "+creatorName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, user := range users {

		if user.Id == creatorId {
			continue
		}

		_, _, err = CreateUserNotification(user.Id, notification.ID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func NotifyEventCreatorOnUserSubscription(creatorId, eventName, userName string) (int, error) {
	notification, _, err := CreateNotification("event", "User "+userName+" subscribed to your event "+eventName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, _, err = CreateUserNotification(creatorId, notification.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func NotifyEventCreatorOnUserUnSubscription(creatorId, eventName, userName string) (int, error) {
	notification, _, err := CreateNotification("event", "User "+userName+" unsubscribed from your event "+eventName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, _, err = CreateUserNotification(creatorId, notification.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func NotifyEventCreatorOnUserComment(creatorId, eventName, userName string) (int, error) {
	notification, _, err := CreateNotification("event", "User "+userName+" commented on your event "+eventName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, _, err = CreateUserNotification(creatorId, notification.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func CreateNotificationSetting(userID string) (int, error) {
	var notificationSetting models.NotificationSetting

	notificationSetting.UserID = userID

	err := db.DB.Create(&notificationSetting).Error

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetNotificationSettingByUserID(userID string) (*types.NotificationSettingsResponse, int, error) {
	var notificationSetting types.NotificationSettingsResponse

	err := db.DB.Model(&models.NotificationSetting{}).Where("user_id = ?", userID).Select("email", "event", "group", "reminder", "id").First(&notificationSetting).Error

	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			CreateNotificationSetting(userID)
			return GetNotificationSettingByUserID(userID)
		}
		return nil, http.StatusInternalServerError, err
	}

	return &notificationSetting, http.StatusOK, nil
}

func UpdateNotificationSetting(userID string, email, event, group, reminder *bool) (int, error) {
	notification, code, err := GetNotificationSettingByUserID(userID)

	if err != nil {
		return code, err
	}

	updateField := make(map[string]interface{})

	if email != nil {
		updateField["Email"] = *email
	}

	if event != nil {
		updateField["Event"] = *event
	}

	if group != nil {
		updateField["Group"] = *group
	}

	if reminder != nil {
		updateField["Reminder"] = *reminder
	}

	err = db.DB.Model(&models.NotificationSetting{}).Where("id", notification.ID).Updates(updateField).Error

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetUsersWithEnabledNotificationSettings(group, event, reminder, email *bool) ([]models.User, int, error) {
	var users []models.User

	if group == nil && event == nil && reminder == nil && email == nil {
		return users, http.StatusConflict, fmt.Errorf("at least one notification setting is needed for this request")
	}

	if group != nil {
		err := db.DB.Table("users").Select("users.*").Joins("LEFT JOIN notification_settings ON notification_settings.user_id = users.id").Where("notification_settings.group = true OR notification_settings.group IS NULL").Find(&users).Error

		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return users, http.StatusOK, nil
	}

	if event != nil {
		err := db.DB.Table("users").Select("users.*").Joins("LEFT JOIN notification_settings ON notification_settings.user_id = users.id").Where("notification_settings.event = true OR notification_settings.event IS NULL").Find(&users).Error

		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return users, http.StatusOK, nil
	}

	if reminder != nil {
		err := db.DB.Table("users").Select("users.*").Joins("LEFT JOIN notification_settings ON notification_settings.user_id = users.id").Where("notification_settings.reminder = true OR notification_settings.reminder IS NULL").Find(&users).Error

		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return users, http.StatusOK, nil
	}

	if email != nil {
		err := db.DB.Table("users").Select("users.*").Joins("LEFT JOIN notification_settings ON notification_settings.user_id = users.id").Where("notification_settings.email = true OR notification_settings.email IS NULL").Find(&users).Error

		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return users, http.StatusOK, nil
	}

	return users, http.StatusOK, nil
}
