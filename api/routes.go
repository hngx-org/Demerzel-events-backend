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

	r.POST("/api/auth/verify", handlers.HandleAuth)

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
	apiRoutes.GET("groups/user", handlers.GetUserGroups)
	apiRoutes.PUT("/groups/:id", handlers.UpdateGroup)
	apiRoutes.GET("/groups/:id", handlers.GetGroupById)
	apiRoutes.DELETE("/groups/:id", handlers.DeleteGroup)

	// User routes
	apiRoutes.GET("/users/current", handlers.GetCurrentUser)
	apiRoutes.GET("/users/:id", handlers.GetUserById)
	apiRoutes.PUT("/users/:id", handlers.UpdateUser)
	apiRoutes.GET("/users", handlers.GetUsers)
	apiRoutes.POST("/users/logout", handlers.LogoutUser)

	// Event Routes
	eventRoutes := apiRoutes.Group("/events")
	eventRoutes.POST("/", handlers.CreateEventHandler)
	eventRoutes.GET("/", handlers.ListEventsHandler)
	eventRoutes.GET("/friends", handlers.ListFriendsEventsHandler)
	eventRoutes.GET("/comments/:event_id", handlers.GetCommentsHandler)
	eventRoutes.POST("join/:event_id", handlers.JoinEventHandler)
	eventRoutes.POST("leave/:event_id", handlers.LeaveEventHandler)
	eventRoutes.GET("/group/:id/", handlers.GetGroupEventsHandler)
	eventRoutes.GET("/subscriptions", handlers.GetUserEventSubscriptions)
	eventRoutes.POST("/:id/subscribe", handlers.SubscribeUserToEvent)
	eventRoutes.POST("/:id/unsubscribe", handlers.UnsubscribeFromEvent)
	eventRoutes.GET("/:event_id", handlers.GetEventHandler)
	// Get events in group returns bad data.

	//comment routes
	commentRoutes := apiRoutes.Group("/comments")
	commentRoutes.POST("/", handlers.CreateComment)
	commentRoutes.PUT("/:comment_id", handlers.UpdateComments)
	commentRoutes.DELETE("/:comment_id", handlers.DeleteComment)

	return r
}
