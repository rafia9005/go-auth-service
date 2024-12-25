package main

import (
	"fmt"
	handler "go-auth-service/internal/handlers"
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

	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/google", handler.AuthGoogle).Methods("GET")
	r.HandleFunc("/auth/google/callback", handler.CallbackAuthGoogle).Methods("GET")
	r.HandleFunc("/auth/github", handler.AuthGithub).Methods("GET")
	r.HandleFunc("/auth/github/callback", handler.CallbackAuthGithub).Methods("GET")
	r.HandleFunc("/verify-token", handler.VerifyToken).Methods("GET")
	r.HandleFunc("/refresh-token", handler.RefreshTokenHandler).Methods("POST")

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
