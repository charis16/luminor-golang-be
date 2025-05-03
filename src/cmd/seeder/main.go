package main

import (
	"fmt"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	config.ConnectDB()
	db := config.DB

	if err := SeedUsers(db); err != nil {
		panic(err)
	}
	fmt.Println("✅ Seeding complete")
}

func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			UUID:        uuid.NewString(),
			Name:        "Admin",
			Email:       "admin@example.com",
			Role:        "admin",
			IsPublished: true,
			Password:    utils.HashPassword("admin123"), // buat fungsi hashPassword
		},
	}

	for _, user := range users {
		var existing models.User
		err := db.Where("email = ?", user.Email).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
