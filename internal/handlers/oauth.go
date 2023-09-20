package handlers

import (
	"demerzel-events/internal/oauth"
	"demerzel-events/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TODO: generate pseudorandoly
var oauth2State = "somerandomoauthstri"

func InitalizeOAuthSignIn(c *gin.Context) {
	// Let generate a random string
	oauth2state := uuid.New().String

	oauthConfig := oauth.OauthConfig()

	// we need a way to store the user to the Oauth2 provider for aunthentication
	c.Session.Set("oauth2_state", oauth2state)
	c.Session.Save()

	c.Redirect(302, oauthConfig.AuthCodeURL(oauth2State))
}

func HandleOAuthCallBack(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	code := c.DefaultQuery("code", "")

	// Let retrieve the stored state from the session
	storedState := c.Session.Get("oauth2_state")

	// let do some checking
	if state != storedState {
		c.Json(http.StatusUnathorized, g.H{
			"message": "Invalid state",
			"status":  "error",
		})
		return
	}

	userInfo, err := services.OAuthCallback(services.NewOAuthCallbackParams(c, code, state))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "An error occurred",
			"status":  "error",
		})
	} else {
		// if there is a successful aunthentication, the stored state should be cleared off from the session
		c.Session.Delete("oauth2_state")
		c.Session.Save()

		c.JSON(200, gin.H{
			"message": "User signed in successfully",
			"status":  "success",
			"data":    userInfo,
		})
	}
}
