package api

import (
	"demerzel-events/pkg/jwt"
	"demerzel-events/pkg/response"
	"demerzel-events/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authentication header is requires")
			c.Abort()
			return
		}

		authToken := authHeader[len("Bearer "):]
		if authToken == "" {
			response.Error(c, http.StatusUnauthorized, "Authentication header is required")
			c.Abort()
			return
		}

		tokenClaims, err := jwt.VerifyToken(authToken, os.Getenv("JWT_SECRET"))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userData := (tokenClaims["data"]).(map[string]interface{})
		if userData["id"] == "" {
			response.Error(c, http.StatusUnauthorized, "The authentication data is incomplete")
			c.Abort()
		}

		user, err := services.GetUserById(userData["id"].(string))
		if err != nil || user == nil {
			response.Error(c, http.StatusUnauthorized, "The authenticated user may have been deleted")
			c.Abort()
		}

		c.Set("user", user)
		c.Next()
	}
}
