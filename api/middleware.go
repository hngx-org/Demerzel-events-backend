package api

import (
	"demerzel-events/pkg/jwt"
	"demerzel-events/pkg/response"
	"github.com/gin-gonic/gin"
	"demerzel-events/services"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.VerifyFromBearer(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userData := (claims["data"]).(map[string]interface{})
		if userData["id"] == "" {
			response.Error(c, http.StatusUnauthorized, "The authentication data is incomplete")
			c.Abort()
			return
		}

		user, err := services.GetUserById(userData["id"].(string))
		if err != nil || user == nil {
			response.Error(c, http.StatusUnauthorized, "The authenticated user may have been deleted")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
