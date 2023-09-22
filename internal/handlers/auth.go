package handlers

import (
	"demerzel-events/dependencies/firebase"
	"demerzel-events/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleAuth(c *gin.Context) {
	var requestBody struct {
		token string
	}

	err := c.BindJSON(&requestBody)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Unable to red request body")
		return
	}

	token, err := firebase.Inner.VerifyIDToken(c, requestBody.token)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Authentication Successful", token)
}
