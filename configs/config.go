package configs

import (
	"demerzel-events/dependencies/firebase"
	"demerzel-events/dependencies/mailersend"
	"demerzel-events/dependencies/mailgun"
	"demerzel-events/internal/db"
	"demerzel-events/pkg/smtp"
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

	// Initialize mailgun
	mailgun.Initialize()

	// Initialize smtp
	smtp.Initialize()

	// Initialize mailersend
	mailersend.Initialize()
}
