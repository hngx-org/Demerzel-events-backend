package api

import (
	"demerzel-events/internal/handlers"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	r.POST("/events", handlers.CreateEventHandler)
	r.GET("/events", handlers.ListEventsHandler)

	// OAuth routes
	oauthRoutes := r.Group("/oauth")

	oauthRoutes.GET("/initialize", handlers.InitalizeOAuthSignIn)
	oauthRoutes.GET("/callback", handlers.HandleOAuthCallBack)

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api")

	// mount the API routes auth middleware
	apiRoutes.Use(AuthMiddleware())

	apiRoutes.POST("/images/upload", handlers.UploadFileHandler)

	// Group routes
	apiRoutes.POST("/groups", handlers.CreateGroup)
	apiRoutes.GET("/groups", handlers.ListGroups)
	apiRoutes.POST("/groups/:id/subscribe", handlers.SubscribeUserToGroup)
	apiRoutes.POST("/groups/:id/unsubscribe", handlers.UnsubscribeFromGroup)
	apiRoutes.PUT("/groups/:id", handlers.UpdateGroup)
	apiRoutes.GET("groups/user", handlers.GetUserGroups)

	// User routes
	apiRoutes.GET("/users/current", handlers.GetCurrentUser)
	apiRoutes.GET("/users/:id", handlers.GetUserById)
	apiRoutes.PUT("/users/:id", handlers.UpdateUser)
	apiRoutes.GET("/users", handlers.GetUsers)

	// Event Routes
	eventRoutes := apiRoutes.Group("/events")
	eventRoutes.GET("", handlers.ListEventsHandler)
	eventRoutes.GET("/:eventid", handlers.GetEventHandler)
	eventRoutes.POST("", handlers.CreateEventHandler)
	eventRoutes.POST("/upload", handlers.UploadFileHandler)
	eventRoutes.GET("/group/:id/", handlers.GroupEventById)

	return r
}
