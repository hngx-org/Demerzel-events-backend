package response

import "github.com/gin-gonic/gin"

type ResponseData = map[string]interface{}

func Success(c *gin.Context, message string, data ResponseData) {
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
