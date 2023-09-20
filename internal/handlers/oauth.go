package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/internal/oauth"
	"demerzel-events/pkg/jwt"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"os"

	"github.com/gin-gonic/gin"
)

// TODO: generate pseudo-randomly
var oauth2State = "somerandomoauthstri"

func InitalizeOAuthSignIn(c *gin.Context) {
	oauthConfig := oauth.OauthConfig()

	responseData := response.ResponseData{"redirectUrl": oauthConfig.AuthCodeURL(oauth2State)}
	response.Success(c, "message", responseData)
}

func HandleOAuthCallBack(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	callbackParams := services.NewOAuthCallbackParams(c, code, state)
	userInfo, err := services.OAuthCallback(callbackParams)
	if err != nil {
		response.Error(c, "An error occurred during authentication")
		return
	}

	user, err := services.GetUserByEmail(userInfo.Email)
	if err != nil {
		response.Error(c, "An error occurred during authentication")
		return
	}

	// Create new user if user does not exist
	if user == nil {
		newUser := models.NewUser(userInfo.Name, userInfo.Email, userInfo.Picture)
		if user, err = services.CreateUser(newUser); err != nil {
			response.Error(c, "An error occurred while creating account")
			return
		}
	}

	data := response.ResponseData{"id": user.Id, "name": user.Name, "email": user.Email}
	token, err := jwt.CreateToken(data, os.Getenv("JWT_SECRET"), 72)
	if err != nil {
		response.Error(c, "An error occurred while generating authentication token")
	}

	response.Success(c, "Authentication successful", response.ResponseData{"token": token})
}
