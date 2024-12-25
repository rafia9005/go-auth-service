package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-auth-service/internal/model/request"
	"go-auth-service/internal/repository"
	service "go-auth-service/internal/services"
	"go-auth-service/pkg/config"
	provider "go-auth-service/pkg/providers"
	"go-auth-service/pkg/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if errValidate := service.ValidateLogin(&loginRequest); errValidate != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Validation failed", "error": errValidate.Error()})
		return
	}

	user, err := service.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
		return
	}

	if !user.Verify {
		utils.RespondJSON(w, http.StatusForbidden, map[string]string{"message": "Account not verified. Please check your email for verification instructions."})
		return
	}

	accessToken, errGenerateAccessToken := service.GenerateAccessToken(*user)
	if errGenerateAccessToken != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating access token"})
		return
	}

	refreshToken, errGenerateRefreshToken := service.GenerateRefreshToken(*user, config.DB)
	if errGenerateRefreshToken != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating refresh token"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":        true,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if errValidate := service.ValidateRegister(&registerRequest); errValidate != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Validation failed", "error": errValidate.Error()})
		return
	}

	result, err := repository.HashAndStoreUser(&registerRequest)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.RespondJSON(w, http.StatusConflict, map[string]string{"message": "Email already in use"})
			return
		}
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to register user"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  result,
		"message": "Registration successful! Please check your email for the verification code",
	})
}

func AuthGoogle(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query().Get("from")
	if form == "" {
		form = "/"
	}
	url := service.GetGoogleAuthURL(form)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackAuthGoogle(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"status": "error", "message": "Authorization code is missing"})
		return
	}

	token, err := provider.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"status": "error", "message": "Failed to exchange authorization code for token"})
		return
	}

	userInfo, err := service.GetGoogleUserInfo(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to get user info: %v", err)})
		return
	}

	email, emailExists := userInfo["email"].(string)
	givenName := userInfo["given_name"].(string)
	familyName := userInfo["family_name"].(string)

	if !emailExists {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "Email is missing from user info"})
		return
	}

	existingUser, err := service.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if saveErr := repository.SaveGoogleUser(givenName, familyName, email); saveErr != nil {
				utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to save new user data: %v", saveErr)})
				return
			}
			existingUser, err = service.GetUserByEmail(email)
			if err != nil {
				utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Failed to fetch the newly created user"})
				return
			}
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to check if user exists: %v", err)})
			return
		}
	}

	if existingUser.Provider != nil && *existingUser.Provider != "google" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": fmt.Sprintf("Your account is already registered with provider '%s'", *existingUser.Provider)})
		return
	}

	accessToken, errGenerateAccessToken := service.GenerateAccessToken(*existingUser)
	if errGenerateAccessToken != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating access token"})
		return
	}

	refreshToken, errGenerateRefreshToken := service.GenerateRefreshToken(*existingUser, config.DB)
	if errGenerateRefreshToken != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating refresh token"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "success",
		"message":       "User authenticated successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func AuthGithub(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query().Get("from")
	if form == "" {
		form = "/"
	}
	url := service.GetGithubAuthUrl(form)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackAuthGithub(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"status": "error", "message": "Authorization code is missing"})
		return
	}

	token, err := provider.GithubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"status": "error", "message": "Failed to exchange authorization code for token"})
		return
	}

	userInfo, err := service.GetGithubUserInfo(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to get user info: %v", err)})
		return
	}

	email, err := service.GetGithubUserPrimaryEmail(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to get user email: %v", err)})
		return
	}

	existingUser, err := service.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var firstName, lastName string
			if name, ok := userInfo["name"].(string); ok {
				nameParts := strings.Fields(name)
				if len(nameParts) > 0 {
					firstName = nameParts[0]
					if len(nameParts) > 1 {
						lastName = strings.Join(nameParts[1:], " ")
					}
				}
			}
			if saveErr := repository.SaveGithubUser(firstName, lastName, email); saveErr != nil {
				utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to save new user data: %v", saveErr)})
				return
			}
			existingUser, err = service.GetUserByEmail(email)
			if err != nil {
				utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Failed to fetch the newly created user"})
				return
			}
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": fmt.Sprintf("Failed to check if user exists: %v", err)})
			return
		}
	}

	if existingUser.Provider != nil && *existingUser.Provider != "github" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"status": "error", "message": fmt.Sprintf("Your account is already registered with provider '%s'", *existingUser.Provider)})
		return
	}

	accessToken, errGenerateAccessToken := service.GenerateAccessToken(*existingUser)
	if errGenerateAccessToken != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating access token"})
		return
	}

	refreshToken, errGenerateRefreshToken := service.GenerateRefreshToken(*existingUser, config.DB)
	if errGenerateRefreshToken != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating refresh token"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "success",
		"message":       "User authenticated successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")
	if token == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"status": "false", "message": "Unauthorized: Token is missing"})
		return
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"status": "false", "message": "Invalid Token"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Token is valid",
		"claims":  claims,
	})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := service.ValidateRefreshToken(request.RefreshToken, config.DB)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, err := service.GenerateAccessToken(*user)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, errGenerateRefreshToken := service.GenerateRefreshToken(*user, config.DB)
	if errGenerateRefreshToken != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}
