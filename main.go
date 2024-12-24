package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("go auth service running on port: ")

	log.Fatal(http.ListenAndServe(":8080", r))
}
