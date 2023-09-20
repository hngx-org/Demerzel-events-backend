package api

import (
	"demerzel-events/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func BuildRoutesHandler() *gin.Engine {
	r := gin.New()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/health", handlers.HealthHandler)

	// OAuth routes
	oauthRoutes := r.Group("/oauth")


	oauthRoutes.GET("/initialize", handlers.InitalizeOAuthSignIn)
	oauthRoutes.GET("/callback", handlers.HandleOAuthCallBack)

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api")

	// User routes
	apiRoutes.GET("/users/:id", handlers.GetUserById)
	apiRoutes.PUT("/users/:id", handlers.UpdateUser)

	return r
}
