package main

import (
	"fmt"
	handler "go-auth-service/internal/handlers"
	"go-auth-service/pkg/config"
	"go-auth-service/pkg/utils"
	"log"
	"net/http"
	"os"

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

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{"message": "Hello, World!"}
		utils.RespondJSON(w, http.StatusOK, payload)
	})

	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/google", handler.AuthGoogle).Methods("GET")
	r.HandleFunc("/auth/google/callback", handler.CallbackAuthGoogle).Methods("GET")
	r.HandleFunc("/auth/github", handler.AuthGithub).Methods("GET")
	r.HandleFunc("/auth/github/callback", handler.CallbackAuthGithub).Methods("GET")
	r.HandleFunc("/verify-token", handler.VerifyToken).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	fmt.Println("go auth service running on port:", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
