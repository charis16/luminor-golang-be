package services

import (
	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/models"
)

func GetAllUsers(page int, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Hitung total data dulu
	if err := config.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data sesuai limit + offset
	offset := (page - 1) * limit
	if err := config.DB.
		Select("uuid", "name", "email", "photo", "description", "phone_number", "url_instagram", "url_tiktok", "url_facebook", "created_at", "updated_at").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func CreateUser(user models.User) models.User {
	config.DB.Create(&user)
	return user
}
