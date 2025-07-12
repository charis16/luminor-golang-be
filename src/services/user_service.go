package services

import (
	"fmt"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
)

type UserInput struct {
	Name         string `form:"name" binding:"required"`
	Email        string `form:"email" binding:"required,email"`
	Role         string `form:"role" binding:"required"`
	Slug         string
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

func GetUserPortfolioBySlug(slug string) (dto.UserPortfolioResponse, error) {
	var user models.User
	if err := config.DB.Where("slug = ?", slug).First(&user).Error; err != nil {
		return dto.UserPortfolioResponse{}, fmt.Errorf("failed to get user by slug: %v", err)
	}

	if user.UUID == "" {
		return dto.UserPortfolioResponse{}, fmt.Errorf("user not found")
	}

	subQuery := config.DB.
		Table("albums").
		Select("category_id").
		Where("user_id = ? AND is_published = ?", user.ID, true)

	var categories []models.Category
	if err := config.DB.
		Table("categories").
		Where("id IN (?) AND is_published = ?", subQuery, true).
		Find(&categories).Error; err != nil {
		return dto.UserPortfolioResponse{}, err
	}

	categoryRes := make([]dto.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		categoryRes = append(categoryRes, dto.CategoryResponse{
			UUID:        c.UUID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
			IsPublished: c.IsPublished,
		})
	}

	response := dto.UserPortfolioResponse{
		User: dto.UserResponse{
			UUID:         user.UUID,
			Name:         user.Name,
			Email:        user.Email,
			Slug:         user.Slug,
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
		},
		Categories: categoryRes, // Replace with the correct field name from dto.UserPortfolioResponse
	}

	return response, nil

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
	tx := config.DB.Begin()
	if tx.Error != nil {
		return models.User{}, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	slug := input.Slug
	if slug == "" {
		slug = utils.GenerateSlug(input.Name)
	} else {
		slug = utils.GenerateSlug(slug)
	}

	var count int64
	if err := tx.Model(&models.User{}).
		Where("slug = ? ", slug).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("failed to check slug uniqueness: %v", err)
	}
	if count > 0 {
		tx.Rollback()
		return models.User{}, fmt.Errorf("slug already exists")
	}

	user := models.User{
		Slug:         slug,
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

	slug := input.Slug
	if slug == "" {
		slug = utils.GenerateSlug(input.Name)
	} else {
		slug = utils.GenerateSlug(slug)
	}

	// Cek apakah slug sudah ada di user lain (selain user ini sendiri)
	var count int64
	if err := tx.Model(&models.User{}).
		Where("slug = ? AND id != ?", slug, user.ID).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("failed to check slug uniqueness: %v", err)
	}
	if count > 0 {
		tx.Rollback()
		return models.User{}, fmt.Errorf("slug already exists")
	}

	// Update field
	user.Slug = slug
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
		if err := utils.DeleteFromR2("users", user.Photo); err != nil {
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
		if len(album.Images) > 0 {
			images = album.Images

			for _, image := range images {
				if image != "" {
					if err := utils.DeleteFromR2("albums", image); err != nil {
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
		if err := utils.DeleteFromR2("users", user.Photo); err != nil {
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

func GetUserOptions() ([]dto.UserResponse, error) {
	var users []models.User
	if err := config.DB.
		Where("is_published = ?", true).
		Where("role != ?", "admin").
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get user options: %v", err)
	}

	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{
			UUID:         user.UUID,
			Name:         user.Name,
			Slug:         user.Slug,
			Photo:        user.Photo,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Role:         user.Role,
			Email:        user.Email,
			PhoneNumber:  user.PhoneNumber,
			URLInstagram: user.URLInstagram,
			URLTikTok:    user.URLTiktok,
			URLFacebook:  user.URLFacebook,
			URLYoutube:   user.URLYoutube,
			IsPublished:  user.IsPublished,
			Description:  user.Description,
		}
	}

	return response, nil
}

func GetTeamMembers() ([]dto.UserResponse, error) {
	var users []models.User
	if err := config.DB.
		Where("is_published = ?", true).
		Where("role != ?", "admin").
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get team members: %v", err)
	}

	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{
			UUID:         user.UUID,
			Name:         user.Name,
			Email:        user.Email,
			Slug:         user.Slug,
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

	return response, nil
}
