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
	r.POST("/api/auth/refresh", AuthMiddleware(), handlers.HandleTokenRefresh)

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
	apiRoutes.POST("/groups/:id/subscribe", handlers.SubscribeUserToGroup)
	apiRoutes.POST("/groups/:id/unsubscribe", handlers.UnsubscribeUserFromGroup)

	// User routes
	apiRoutes.GET("/users/current", handlers.GetCurrentUser)
	apiRoutes.PUT("/users", handlers.UpdateUser)
	apiRoutes.GET("/users", handlers.GetUsers)
	apiRoutes.POST("/users/logout", handlers.LogoutUser)

	// Event Routes
	apiRoutes.POST("/events", handlers.CreateEventHandler)
	apiRoutes.GET("/events", handlers.ListEventsHandler)
	apiRoutes.GET("/events/upcoming", handlers.ListUpcomingEventsHandler)
	apiRoutes.GET("/events/friends", handlers.ListFriendsEventsHandler)
	apiRoutes.GET("/events/comments/:event_id", handlers.GetCommentsHandler)
	apiRoutes.GET("/events/group/:id", handlers.GetGroupEventsHandler)
	apiRoutes.GET("/events/subscriptions", handlers.GetUserEventSubscriptions)
	apiRoutes.POST("/events/:id/subscribe", handlers.SubscribeUserToEvent)
	apiRoutes.POST("/events/:id/unsubscribe", handlers.UnsubscribeFromEvent)
	apiRoutes.GET("/events/:event_id/attendees", handlers.GetEventAttendees)
	apiRoutes.GET("/events/:event_id", handlers.GetEventHandler)
	apiRoutes.PUT("/events/:event_id", handlers.UpdateEvent)
	apiRoutes.DELETE("/events/:event_id", handlers.DeleteEventHandler)

	apiRoutes.POST("/events/:id/reactions", handlers.ToggleEventReaction)
	apiRoutes.GET("/events/:event_id/reactions", handlers.GetEventReactions)
	apiRoutes.DELETE("/events/:event_id/reactions/:reaction_id", handlers.RemoveReaction)

	// comment routes
	apiRoutes.POST("/comments", handlers.CreateComment)
	apiRoutes.PUT("/comments/:comment_id", handlers.UpdateComments)
	apiRoutes.DELETE("/comments/:comment_id", handlers.DeleteComment)

	// Notification routes
	apiRoutes.GET("/notifications/all", handlers.ListNotifications)
	apiRoutes.GET("/notifications", handlers.GetUserNotifications)
	apiRoutes.GET("/notifications/settings", handlers.GetNotificationSettings)
	apiRoutes.PUT("/notifications/settings", handlers.UpdateNotificationSettings)
	apiRoutes.PUT("/notifications/:id", handlers.UpdateUserNotification)
	apiRoutes.PUT("/notifications", handlers.UpdateUserNotifications)

	// Tag routes
	apiRoutes.GET("/tags", handlers.GetTags)

	return r
}
