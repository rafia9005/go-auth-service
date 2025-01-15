package handlers

import (
	"go-auth-service/internal/model/entity"
	"go-auth-service/pkg/config"
	"go-auth-service/pkg/utils"
	"net/http"
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
		"contact": map[string]interface{}{
			"phone": user.Contacts[0].Phone,
		},
		"refresh_token": user.RefreshTokens[0].Token,
		"verify":        user.Verify,
		"provider":      *user.Provider,
    "address": map[string]interface{}{
      "city": "Malang",
      "country": "Indonesian",
    },
		"created_at":    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.RespondJSON(w, http.StatusOK, filter)
}
