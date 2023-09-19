package handlers

import (
	"demerzel-events/internal/oauth"
	"demerzel-events/services"
	"github.com/gin-gonic/gin"
)

// TODO: generate pseudorandoly
var oauth2State = "somerandomoauthstri"

func InitalizeOAuthSignIn(c *gin.Context) {
	oauthConfig := oauth.OauthConfig()

	c.Redirect(302, oauthConfig.AuthCodeURL(oauth2State))
}

func HandleOAuthCallBack(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	code := c.DefaultQuery("code", "")

	userInfo, err := services.OAuthCallback(services.NewOAuthCallbackParams(c, code, state))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "An error occurred",
			"status":  "error",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "User signed in successfully",
			"status":  "success",
			"data":    userInfo,
		})
	}
}
