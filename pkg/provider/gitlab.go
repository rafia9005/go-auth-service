package provider

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/gitlab"
)

var GitlabOauthConfig *oauth2.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	GithubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITLAB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITLAB_CLIENT_SECRET_KEY"),
		RedirectURL:  "http://localhost:3000/api/v1/auth/gitlab/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     gitlab.Endpoint,
	}
}
