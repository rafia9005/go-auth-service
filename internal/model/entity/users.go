package entity

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID            uint           `gorm:"primaryKey"`
	Name          string         `json:"name"`
	FirstName     string         `json:"first_name"`
	LastName      *string        `json:"last_name"`
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	Role          string         `json:"role" gorm:"type:enum('admin','member')"`
	Verify        bool           `json:"verify"`
	Provider      *string        `json:"provider" gorm:"type:enum('default', 'google', 'github');default:'default'"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`
  Contacts      []Contacts     `gorm:"foreignKey:UserID"`
}
