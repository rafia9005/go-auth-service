package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-auth-service/internal/middleware"
	"go-auth-service/internal/model/entity"
	"go-auth-service/internal/model/request"
	"go-auth-service/pkg/config"
	provider "go-auth-service/pkg/providers"
	"go-auth-service/pkg/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
)

// validate login req
func ValidateLogin(loginRequest *request.LoginRequest) error {
	validate := validator.New()
	return validate.Struct(loginRequest)
}

// get user by emails
func GetUserByEmail(email string) (*entity.Users, error) {
	var user entity.Users
	err := config.DB.First(&user, "email = ?", email).Error
	return &user, err
}

// generate jwt token
func GenerateJWTToken(user *entity.Users) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
		"role":  "member",
	}

	if user.Role == "admin" {
		claims["role"] = "admin"
	}

	return utils.GenerateToken(&claims)
}

// validate register
func ValidateRegister(registerRequest *request.RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(registerRequest)
}

// auth users
func AuthenticateUser(email, password string) (*entity.Users, error) {
	var user entity.Users
	err := config.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	if !middleware.CheckPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

// google auth callback url
func GetGoogleAuthURL(redirectURI string) string {
	return provider.GoogleOauthConfig.AuthCodeURL(redirectURI)
}

// github auth callback url
func GetGithubAuthUrl(redirectURI string) string {
	return provider.GithubOauthConfig.AuthCodeURL(redirectURI)
}

// google users info
func GetGoogleUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := provider.GoogleOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return userInfo, nil
}


// github users info
func GetGithubUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := provider.GithubOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return userInfo, nil
}

// github emails users info
func GetGithubUserPrimaryEmail(token *oauth2.Token) (string, error) {
	client := provider.GithubOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user emails: status code %d", resp.StatusCode)
	}

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode user emails: %w", err)
	}

	for _, e := range emails {
		if primary, ok := e["primary"].(bool); ok && primary {
			if email, ok := e["email"].(string); ok {
				return email, nil
			}
		}
	}

	return "", fmt.Errorf("no primary email found")
}
