package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/internal/oauth"
	"demerzel-events/pkg/jwt"
	"demerzel-events/pkg/response"
	"demerzel-events/pkg/types"
	"demerzel-events/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// TODO: generate this pseudo-randomly
var oauth2State = "somerandomoauthstri"

func InitalizeOAuthSignIn(c *gin.Context) {
	oauthConfig := oauth.OauthConfig()

	data := types.ResponseMap{"redirectUrl": oauthConfig.AuthCodeURL(oauth2State)}
	response.Success(c, http.StatusOK, "Authentication initialized", data)
}

func HandleOAuthCallBack(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	callbackParams := services.NewOAuthCallbackParams(c, code, state)
	userInfo, err := services.OAuthCallback(callbackParams)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "An error occurred during authentication")
		return
	}

	user, err := services.GetUserByEmail(userInfo.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "An error occurred during authentication")
		return
	}

	// Create new user if user does not exist
	if user == nil {
		newUser := models.NewUser(userInfo.Name, userInfo.Email, userInfo.Picture)
		if user, err = services.CreateUser(newUser); err != nil {
			response.Error(c, http.StatusInternalServerError, "An error occurred while creating account")
			return
		}
	}

	data := types.ResponseMap{"id": user.Id, "name": user.Name, "email": user.Email}
	token, err := jwt.CreateToken(data, os.Getenv("JWT_SECRET"), 72)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "An error occurred while generating authentication token")
	}

	refreshToken, err := jwt.CreateToken(types.ResponseMap{"id": user.Id}, os.Getenv("JWT_SECRET"), 730)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "An error occurred while generating authentication token")
		return
	}

	response.Success(c, http.StatusOK, "Authentication successful", types.ResponseMap{"token": token, "refreshToken": refreshToken})
}
