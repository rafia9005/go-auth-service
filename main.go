package main

import (
	"fmt"
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

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{"message": "Hello, World!"}
		utils.RespondJSON(w, http.StatusOK, payload)
	})

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	fmt.Println("go auth service running on port:", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
