package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

func OauthConfig() oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}
