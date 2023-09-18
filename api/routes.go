package api

import (
	"Demerzel-Events/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Application struct{}

func BuildRoutes() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/health", handlers.HealthHandler)

	return r
}
