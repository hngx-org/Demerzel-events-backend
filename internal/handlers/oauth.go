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
		Id: userInfo.Sub,
        Email: userInfo.Email,
        Name:  userInfo.Name,
		Avatar: userInfo.Picture,
    }
	if existingUser != nil {
// Set the ID to update the existing user.
	newUser.Id = existingUser.Id

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
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})

	}
	user, err := services.GetUserFromDB(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})
	}

	c.JSON(http.StatusOK, user)

}

func HandleUpdateUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := services.GetUserFromDB(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
        return
    }

    if user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
	var updatedUser models.UpdateUserStruct
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Update the user's data
    user.Email = updatedUser.Email
    user.Name = updatedUser.Name

    // Save the updated user
    if err := services.UpdateUserByID(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
        return
    }

    // Return the updated user as JSON
    c.JSON(http.StatusOK, user)

}
