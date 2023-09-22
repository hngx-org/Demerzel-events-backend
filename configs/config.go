package configs

import (
	"demerzel-events/internal/db"
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
}
