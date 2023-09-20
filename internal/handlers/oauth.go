package handlers

import (
	"demerzel-events/internal/models"
	"demerzel-events/internal/oauth"
	"demerzel-events/pkg/jwt"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"
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

        if err := services.UpdateUserByID(newUser); err != nil {
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

	//user, err := services.GetUserByEmail(userInfo.Email)

	if err != nil {
		response.Error(c, "An error occurred during authentication")
		return
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


	// Create new user if user does not exist
	if user == nil {
		newUser := models.NewUser(user.Name, user.Email, user.Avatar)
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
