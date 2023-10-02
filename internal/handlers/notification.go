package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/pkg/types"
	"demerzel-events/services"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetUserNotifications(ctx *gin.Context) {

	rawUser, exists := ctx.Get("user")

	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(ctx, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	notifications, status, err := services.ListUserNotifications(user.Id)

	if err != nil {
		response.Error(ctx, status, err.Error())
		return
	}

	if notifications == nil {
		notifications = []types.UserNotificationResponse{}
	}

	response.Success(ctx, status, "List of notifications", map[string]interface{}{"notifications": notifications})
	
}

func UpdateUserNotification(ctx *gin.Context) {
	rawUser, exists := ctx.Get("user")
	id := ctx.Param("id")

	if !exists {
		response.Error(ctx, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)

	if !ok {
		response.Error(ctx, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	var requestBody struct {
		Read bool `json:"read" validate:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	getNotification, status, err := services.GetUserNotificationByID(id)

	if err != nil {
		response.Error(ctx, status, err.Error())
		return
	}

	if getNotification.UserID != user.Id {
		response.Error(ctx, http.StatusForbidden, "You are not allowed to update this notification")
		return
	}

	 code, err := services.UpdateUserNotification(id, requestBody.Read)

	if err != nil {
		response.Error(ctx, code, err.Error())
		return
	}


	notificationResponse, code, err := services.GetUserNotificationByID(id)

	if err != nil {
		response.Error(ctx, code, err.Error())
		return
	}

	

	response.Success(ctx, code, "Notification updated successfully", map[string]interface{}{"notification": notificationResponse})

}

func ListNotifications(ctx *gin.Context) {
	notifications, status, err := services.ListNotifications()

	if err != nil {
		response.Error(ctx, status, err.Error())
		return
	}

	if notifications == nil {
		notifications = []models.Notification{}
	}

	response.Success(ctx, status, "List of notifications", map[string]interface{}{"notifications": notifications})
}
