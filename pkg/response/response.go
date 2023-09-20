package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"status":  "error",
		"message": message,
	})
}

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}