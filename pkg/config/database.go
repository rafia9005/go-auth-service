package config

import (
	"fmt"
	"go-auth-service/internal/model/entity"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() error {
	dsn := os.Getenv("MYSQL")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect database!")
		return err
	}

	DB = db

	DB.AutoMigrate(&entity.Users{}, &entity.RefreshToken{})

	return nil
}
