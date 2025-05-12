package services

import (
	"encoding/json"
	"fmt"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
)

type UserInput struct {
	Name         string
	Email        string
	Role         string
	Description  string
	PhotoURL     string
	Password     string
	URLInstagram string
	URLTikTok    string
	URLFacebook  string
	URLYoutube   string
	PhoneNumber  string
	CanLogin     bool
	IsPublished  bool // tetap string kalau dari form
}

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
			URLTikTok:    user.URLTiktok,
			URLFacebook:  user.URLFacebook,
			URLYoutube:   user.URLYoutube,
			IsPublished:  user.IsPublished,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		}
	}

	return response, total, nil
}

func GetUserByUUID(uuid string) (models.User, error) {
	var user models.User
	if err := config.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

func CreateUser(input UserInput) (models.User, error) {
	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		Role:         input.Role,
		Description:  input.Description,
		URLInstagram: input.URLInstagram,
		URLTiktok:    input.URLTikTok,
		URLFacebook:  input.URLFacebook,
		URLYoutube:   input.URLYoutube,
		PhoneNumber:  input.PhoneNumber,
		IsPublished:  input.IsPublished,
	}

	if input.Password != "" {
		user.Password = utils.HashPassword(input.Password)
	}
	if input.PhotoURL != "" {
		user.Photo = input.PhotoURL
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		return models.User{}, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("failed to save user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return models.User{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return user, nil
}

func UpdateUser(uuid string, input UserInput) (models.User, error) {
	var user models.User

	tx := config.DB.Begin()
	if tx.Error != nil {
		return models.User{}, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Cari user
	if err := tx.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("failed to find user: %v", err)
	}

	// Update field
	user.Name = input.Name
	user.Email = input.Email
	user.Role = input.Role
	user.Description = input.Description
	user.Photo = input.PhotoURL
	if input.Password != "" && input.CanLogin {
		user.Password = utils.HashPassword(input.Password)
	}

	if !input.CanLogin {
		user.Password = ""
	}
	user.URLInstagram = input.URLInstagram
	user.URLTiktok = input.URLTikTok
	user.URLFacebook = input.URLFacebook
	user.URLYoutube = input.URLYoutube
	user.PhoneNumber = input.PhoneNumber
	user.IsPublished = input.IsPublished

	// Simpan perubahan
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("failed to update user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return models.User{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return user, nil
}

func DeleteUser(uuid string) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	var user models.User
	if err := tx.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get user: %v", err)
	}

	// === Step 1: Delete user photo from MinIO bucket "users"
	if user.Photo != "" {
		if err := utils.DeleteFromMinio("users", user.Photo); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete user photo: %v", err)
		}
	}

	// === Step 2: Get all albums owned by user
	var albums []models.Album
	if err := tx.Where("user_id = ?", user.ID).Find(&albums).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get albums: %v", err)
	}

	for _, album := range albums {
		var images []string
		if album.Images != "" {
			err := json.Unmarshal([]byte(album.Images), &images)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to parse album images: %v", err)
			}

			for _, image := range images {
				if image != "" {
					if err := utils.DeleteFromMinio("albums", image); err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to delete album image: %v", err)
					}
				}
			}
		}
	}
	// === Step 3: Delete albums from DB
	if err := tx.Where("user_id = ?", user.ID).Delete(&models.Album{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete albums: %v", err)
	}

	// === Step 4: Delete user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// === Step 5: Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func DeleteImageUser(uuid string) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	var user models.User
	if err := tx.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.Photo != "" {
		if err := utils.DeleteFromMinio("users", user.Photo); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete user photo: %v", err)
		}
	}

	if err := tx.Model(&user).Updates(map[string]interface{}{
		"photo": nil,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user photo: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
