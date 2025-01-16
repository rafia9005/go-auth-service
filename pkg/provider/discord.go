package provider

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
  "golang.org/x/oauth/discord"
)

var DiscordOauthConfig *oauth2.Config

func init() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  DiscordOauthConfig = &oauth2.Config{
    ClientID: os.Getenv("DISCORD_CLIENT_ID"),
    ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET_ID"),
    RedirectURL: "http://localhost:3000/api/v1/auth/discord/callback",
    Scopes: []string{"email", "identify"},
    Endpoint: discord.Endpoint,
  }
}
