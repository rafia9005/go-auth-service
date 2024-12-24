package repository

import (
	"fmt"
	"go-auth-service/internal/middleware"
	"go-auth-service/internal/model/entity"
	"go-auth-service/internal/model/request"
	"go-auth-service/pkg/config"
)

// update users data
func UpdateUser(user *entity.Users) error {
	return config.DB.Save(user).Error
}

// hash and store users
func HashAndStoreUser(registerRequest *request.RegisterRequest) (string, error) {
	var existingUser entity.Users
	if err := config.DB.First(&existingUser, "email = ?", registerRequest.Email).Error; err == nil {
		return "", fmt.Errorf("user with email %s already exists", registerRequest.Email)
	}

	hashedPassword, err := middleware.HashPassword(registerRequest.Password)
	if err != nil {
		return "", err
	}

	newUser := entity.Users{
		Name:      fmt.Sprintf("%s %s", registerRequest.FirstName, registerRequest.LastName),
		FirstName: registerRequest.FirstName,
		LastName:  &registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  hashedPassword,
		Role:      "member",
		Verify:    true,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("User %s registered successfully", newUser.Email), nil
}

// save google users
func SaveGoogleUser(firstName, lastName, email string) error {
	provider := "google"
	newUser := entity.Users{
		Name:      fmt.Sprintf("%s %s", firstName, lastName),
		FirstName: firstName,
		LastName:  &lastName,
		Email:     email,
		Role:      "member",
		Verify:    true,
		Provider:  &provider,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return err
	}
	return nil
}

// save github users
func SaveGithubUser(fristName, lastName, email string) error {
	provider := "github"
	newUser := entity.Users{
		Name:      fmt.Sprintf("%s %s", fristName, lastName),
		FirstName: fristName,
		LastName:  &lastName,
		Email:     email,
		Role:      "member",
		Verify:    true,
		Provider:  &provider,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return err
	}

	return nil
}
