package api

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"demerzel-events/internal/handlers"
	"demerzel-events/services"
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

	// init group service
	groupService := services.NewGroupService()
	// create group hanlder
	groupHandler := handlers.Group{
		Service: groupService,
	}
	r.GET("/groups", groupHandler.List)

	return r
}
