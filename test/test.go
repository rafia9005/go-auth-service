package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
)

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "logged in successfully"}`))
}

func TestLogin(t *testing.T) {
	t.Parallel()

	body := strings.NewReader(`{"email": "rafia9005@gmail.com", "password": "admin123"}`)
	req, err := http.NewRequest("POST", "/api/v1/auth/login", body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/auth/login", LoginRequest).Methods("POST")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := `{"message": "logged in successfully"}`
	assert.Equal(t, expectedResponse, rr.Body.String())
}
