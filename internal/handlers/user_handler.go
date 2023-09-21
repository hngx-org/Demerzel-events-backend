package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "demerzel-events/internal/models/user"
)

func CreateUser(c *gin.Context) {
     var newUser models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

   c.JSON(http.StatusCreated, user)
}

func GetUserByID(c *gin.Context) {
    userID := c.Param("id") // If ID is the URL parameter for the user ID
    var existingUser models.User
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}


func UpdateUser(c *gin.Context) {
    userID := c.Param("id") // If ID is the URL parameter for the user ID
    var existingUser models.User
    if err := db.Where("id = ?", userID).First(&existingUser).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if !isUserAuthorized(c) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var updatedUser models.User
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    existingUser.Username = updatedUser.Username
    existingUser.Email = updatedUser.Email
    // Update other fields as needed


    if err := db.Save(&existingUser).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }
    c.JSON(http.StatusOK, existingUser)
}

func DeleteUser(c *gin.Context) {
    userID := c.Param("id") // If ID is the URL parameter for the user ID

    var userToDelete models.User
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if !isUserAuthorized(c) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    if err := db.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}


func isUserAuthorized(c *gin.Context) bool {
    
    //authorization logic here.
 
    return true   // Return true if authorized; otherwise, return false.
 
}
