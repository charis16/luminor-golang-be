package services

import (
	"fmt"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
)

func GetAllUsers(page int, limit int, search string) ([]dto.UserResponse, int64, error) {
	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{})

	fmt.Printf("Search term: %s\n", search)
	// Apply search filter if search term is provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Select("uuid", "name", "email", "photo", "description", "phone_number", "url_instagram", "url_tiktok", "url_facebook", "created_at", "updated_at", "role").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// Mapping ke response DTO
	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{
			UUID:         user.UUID,
			Name:         user.Name,
			Email:        user.Email,
			Photo:        user.Photo,
			Description:  user.Description,
			Role:         user.Role,
			PhoneNumber:  user.PhoneNumber,
			URLInstagram: user.URLInstagram,
			URLTiktok:    user.URLTiktok,
			URLFacebook:  user.URLFacebook,
			URLYoutube:   user.URLYoutube,
			IsPublished:  user.IsPublished,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		}
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
