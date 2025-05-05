package services

import (
	"fmt"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
)

func GetAllUsers(page int, limit int) ([]dto.UserResponse, int64, error) {
	var users []models.User
	var total int64

	if err := config.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := config.DB.
		Select("uuid", "name", "email", "photo", "description", "phone_number", "url_instagram", "url_tiktok", "url_facebook", "created_at", "updated_at").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// Mapping ke response DTO
	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{User: user}
	}

	return response, total, nil
}

func CreateUser(name, email, role, description string, photoURL string, password string, urlInstagram string, url_tikTok string, urlFacebook string, urlYoutube string, phoneNumber string, isPublished string) (models.User, error) {

	user := models.User{
		Name:         name,
		Email:        email,
		Role:         role,
		Description:  description,
		Photo:        photoURL,
		Password:     utils.HashPassword(password),
		URLInstagram: urlInstagram,
		URLTiktok:    url_tikTok,
		URLFacebook:  urlFacebook,
		URLYoutube:   urlYoutube,
		PhoneNumber:  phoneNumber,
		IsPublished:  isPublished == "true",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("failed to save user: %v", err)
	}

	return user, nil
}
