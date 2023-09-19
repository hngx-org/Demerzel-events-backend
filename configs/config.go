package configs

import (
	"demerzel-events/internal/db"
	"github.com/joho/godotenv"
)

func Load() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Setup database connection
	db.SetupDB()
}
