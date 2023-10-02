package services

import (
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/types"
	"net/http"
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

func SendNewGroupNotificationToAllUsers(groupName, creatorName, creatorId string) (int, error) {
	var users []models.User

	err := db.DB.Find(&users).Error

	if err != nil {
		return http.StatusInternalServerError, err
	}

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
