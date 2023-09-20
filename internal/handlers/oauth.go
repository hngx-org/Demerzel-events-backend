package handlers

import (
	"demerzel-events/internal/oauth"
	"demerzel-events/services"
	"net/http"

	"demerzel-events/internal/models"

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
	existingUser, err := services.GetUserByEmail(userInfo.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing user"})
        return
    }
	newUser := &models.UserData{
        Email: userInfo.Email,
        Name:  userInfo.Name,
    }
	if existingUser != nil {
// Set the ID to update the existing user.
        if err := services.UpdateUserByEmail(newUser); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
            return
        }
    } else {
        if err := services.CreateUser(newUser); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
            return
        }
    }
	c.JSON(http.StatusOK, gin.H{"message": "User authenticated and data saved"})
}

func HandleGetUserID(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})

	}
	var user models.User
	user, err := services.GetUserFromDB(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})
	}

	c.JSON(http.StatusOK, user)

}

func HandleUpdateUser(c *gin.Context) {
	var UpdateUserStruct models.UpdateUserStruct
	if err := c.ShouldBindJSON(&UpdateUserStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
	}
	if err := services.UpdateUser(UpdateUserStruct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error Updating the User",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Have Been Updated Successfully"})
}
