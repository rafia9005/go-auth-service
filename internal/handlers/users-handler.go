package handlers

import (
	"encoding/json"
	"go-auth-service/internal/model/entity"
	"go-auth-service/internal/model/request"
	"go-auth-service/pkg/config"
	"go-auth-service/pkg/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	if token == "" {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized",
		})
		return
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized: invalid user_id claim",
		})
		return
	}

	var user entity.Users
	if err := config.DB.Preload("Contacts").Preload("RefreshTokens").First(&user, uint(userID)).Error; err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "User not found",
		})
		return
	}

	filter := map[string]interface{}{
		"name":       user.Name,
		"first_name": user.FirstName,
		"last_name":  *user.LastName,
		"email":      user.Email,
		"role":       user.Role,
		"verify":     user.Verify,
		"provider":   *user.Provider,
		"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at": user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.RespondJSON(w, http.StatusOK, filter)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")
	claims, err := utils.DecodeToken(token)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized",
		})
		return
	}

	userID := claims["user_id"].(float64)

	var updateRequest request.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(updateRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	var user entity.Users
	if err := config.DB.First(&user, uint(userID)).Error; err != nil {
		utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
			"message": "User not found",
		})
		return
	}

	// validate if nil
	if updateRequest.FirstName != "" {
		user.FirstName = updateRequest.FirstName
	}

	if updateRequest.LastName != "" {
		user.LastName = &updateRequest.LastName
	}

	if updateRequest.Name != "" {
		user.Name = updateRequest.Name
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to update user",
		})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
	})
}
