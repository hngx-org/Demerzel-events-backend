package handlers

import (
	"demerzel-events/dependencies/firebase"
	"demerzel-events/internal/models"
	"demerzel-events/pkg/jwt"
	"demerzel-events/pkg/response"
	"demerzel-events/pkg/types"
	"demerzel-events/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleAuth(c *gin.Context) {
	var requestBody struct {
		Token string `json:"token"`
	}

	err := c.BindJSON(&requestBody)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Unable to red request body")
		return
	}

	claims, err := firebase.Inner.VerifyIDToken(requestBody.Token)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	email := claims["email"].(string)
	user, err := services.GetUserByEmail(email)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if user == nil {
		newUser := models.NewUser(claims["name"].(string), email, claims["picture"].(string))
		user, err = services.CreateUser(newUser)

		if err != nil {
			response.Error(c, http.StatusInternalServerError, "An error occurred while creating account")
			return
		}
	}

	token, err := jwt.CreateToken(types.ResponseMap{"id": user.Id}, os.Getenv("JWT_SECRET"), 72)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "An error occurred while generating authentication token")
	}

	response.Success(c, http.StatusOK, "Authentication Successful", types.ResponseMap{"token": token})
}
