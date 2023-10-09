package configs

import (
	"demerzel-events/dependencies/firebase"
	"demerzel-events/internal/db"
	"demerzel-events/internal/models"
	"demerzel-events/services"
	"fmt"

	"github.com/joho/godotenv"
)

func Load() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error: ccannot find .env file in the project root")
	}

	// Setup database connection
	db.SetupDB()

	// Initialize firebase
	firebase.Initialize()

	// Check if tags table hasn't been populated already(if it is empty)
	var tagsCount int64
	db.DB.Model(&models.Tag{}).Count(&tagsCount)

	if tagsCount == 0 {
		services.PrepopulateTags()
	}
}
