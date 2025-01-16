package main

import (
	"fmt"
	handler "go-auth-service/internal/handlers"
	"go-auth-service/internal/middleware"
	"go-auth-service/pkg/config"
	"go-auth-service/pkg/logs"
	"go-auth-service/pkg/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load the env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect database
	config.DBConnect()

	r := mux.NewRouter()

	r.Use(logs.Logging)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{
			"status":  "Server is running smoothly 🚀",
			"version": "1.0.0",
			"message": "Welcome to our awesome API! 🎉",
		}
		utils.RespondJSON(w, http.StatusOK, payload)
	})

	router := r.PathPrefix("/api/v1").Subrouter()

	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/google", handler.AuthGoogle).Methods("GET")
	router.HandleFunc("/auth/google/callback", handler.CallbackAuthGoogle).Methods("GET")
	router.HandleFunc("/auth/github", handler.AuthGithub).Methods("GET")
	router.HandleFunc("/auth/github/callback", handler.CallbackAuthGithub).Methods("GET")
  router.HandleFunc("/auth/gitlab", handler.AuthGithub).Methods("GET")
  router.HandleFunc("/auth/gitlab/callback", handler.CallbackAuthGithub).Methods("GET")
	router.HandleFunc("/auth/verify-token", handler.VerifyToken).Methods("GET")
	router.HandleFunc("/auth/refresh-token", handler.RefreshTokenHandler).Methods("POST")

  protected := router.NewRoute().Subrouter()
  protected.Use(middleware.Auth)
  protected.HandleFunc("/profile", handler.Profile).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	logServiceStart(port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func logServiceStart(port string) {
	startTime := time.Now().Format(time.RFC1123)
	message := fmt.Sprintf("🚀 Service running on http://localhost:%s | Started at: %s", port, startTime)
	log.Println(message)
}
