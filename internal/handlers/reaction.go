package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ToggleEventReaction(c *gin.Context) {
	eventId := c.Param("event_id")

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	var requestBody struct {
		Reaction string `json:"reaction"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body format")
		return
	}

	reaction, err := services.GetReactionForEvent(user.Id, eventId, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if reaction == nil {
		err = services.AddReaction(user, eventId, requestBody.Reaction)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		err = services.UpdateReaction(user, reaction.Id, requestBody.Reaction)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("Reaction \"%s\" added successfully", requestBody.Reaction), nil)
}

func GetEventReactions(c *gin.Context) {
	eventId := c.Param("event_id")

	reactions, err := services.GetAllReactionsForEvent(eventId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Reactions retrieved successfully", reactions)
}

func RemoveReaction(c *gin.Context) {
	eventId := c.Param("event_id")
	reactionId := c.Param("reaction_id")

	rawUser, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "Unable to read user from context")
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid context user type")
		return
	}

	reaction, err := services.GetReactionById(reactionId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if reaction == nil {
		response.Error(c, http.StatusNotFound, "reaction does not exist")
		return
	}

	if reaction.EventId != eventId {
		response.Error(c, http.StatusConflict, "reaction and event mismatch")
		return
	}

	err = services.RemoveReaction(user, reaction.Id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Reaction removed successfully", nil)
}
