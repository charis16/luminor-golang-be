package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
)

type CategoryInput struct {
	Name        string `form:"name" validate:"required"`
	Slug        string `form:"slug" validate:"required"`
	Description string `form:"description" validate:"required"`
	IsPublished string `form:"is_published" validate:"required"`
	PhotoUrl    string `form:"-"` // handled manually
}

func GetPublishedCategories() ([]dto.CategoryResponse, error) {
	var categories []models.Category

	if err := config.DB.Where("is_published = ?", true).
		Select("uuid", "name", "is_published", "slug", "description", "photo_url", "created_at", "updated_at").
		Order("created_at DESC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	// Mapping ke response DTO
	response := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = dto.CategoryResponse{
			UUID:        category.UUID,
			Name:        category.Name,
			Description: category.Description,
			Slug:        category.Slug,
			PhotoUrl:    category.PhotoURL,
			IsPublished: category.IsPublished,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		}
	}

	return response, nil
}

func GetAllCategories(page int, limit int, search string) ([]dto.CategoryResponse, int64, error) {
	var categories []models.Category
	var total int64

	query := config.DB.Model(&models.Category{})

	// Apply search filter if search term is provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Select("uuid", "name", "is_published", "created_at", "updated_at").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	// Mapping ke response DTO
	response := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = dto.CategoryResponse{
			UUID:        category.UUID,
			Name:        category.Name,
			IsPublished: category.IsPublished,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		}
	}

	return response, total, nil
}

func CreateCategory(input CategoryInput) (*models.Category, error) {
	tx := config.DB.Begin() // Mulai transaksi

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	if input.Slug != "" {
		var count int64
		if err := tx.Model(&models.User{}).
			Where("slug = ? ", strings.ReplaceAll(input.Slug, " ", "-")).
			Count(&count).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to check slug uniqueness: %v", err)
		}
		if count > 0 {
			tx.Rollback()
			return nil, fmt.Errorf("slug already exists")
		}
	}

	category := models.Category{
		Name:        input.Name,
		IsPublished: input.IsPublished == "1",
		Description: input.Description,
		Slug:        strings.ReplaceAll(input.Slug, " ", "-"),
		PhotoURL:    input.PhotoUrl,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&category).Error; err != nil {
		tx.Rollback() // rollback jika error
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err // commit gagal
	}

	return &category, nil
}

func GetCategoryByUUID(uuid string) (models.Category, error) {
	var category models.Category
	if err := config.DB.Where("uuid = ?", uuid).First(&category).Error; err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func UpdateCategory(uuid string, input CategoryInput) (models.Category, error) {
	tx := config.DB.Begin()

	if input.Slug != "" {
		var count int64
		if err := tx.Model(&models.Category{}).
			Where("slug = ? AND uuid != ?", strings.ReplaceAll(input.Slug, " ", "-"), uuid).
			Count(&count).Error; err != nil {
			tx.Rollback()
			return models.Category{}, fmt.Errorf("failed to check slug uniqueness: %v", err)
		}
		if count > 0 {
			tx.Rollback()
			return models.Category{}, fmt.Errorf("slug already exists")
		}
	}

	var category models.Category

	if err := tx.Where("uuid = ?", uuid).First(&category).Error; err != nil {
		tx.Rollback()
		return models.Category{}, err
	}

	category.Name = input.Name
	category.IsPublished = input.IsPublished == "1"
	category.Description = input.Description
	category.Slug = strings.ReplaceAll(input.Slug, " ", "-")

	if input.PhotoUrl != "" && input.PhotoUrl != "undefined" {
		category.PhotoURL = input.PhotoUrl
	}

	category.UpdatedAt = time.Now()

	if err := tx.Save(&category).Error; err != nil {
		tx.Rollback()
		return models.Category{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func DeleteCategory(uuid string) error {
	tx := config.DB.Begin()
	var category models.Category

	// Cari category berdasarkan UUID
	if err := tx.Where("uuid = ?", uuid).First(&category).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Ambil semua album terkait category ini
	var albums []models.Album
	if err := tx.Where("category_id = ?", category.ID).Find(&albums).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus image album dari MinIO

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

	// Hapus album-album terkait
	if err := tx.Where("category_id = ?", category.ID).Delete(&models.Album{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus category-nya
	if err := tx.Delete(&category).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func DeleteImageCategory(uuid string) error {
	tx := config.DB.Begin()
	var category models.Category

	// Cari category berdasarkan UUID
	if err := tx.Where("uuid = ?", uuid).First(&category).Error; err != nil {
		tx.Rollback()
		return err
	}

	if category.PhotoURL != "" {
		if err := utils.DeleteFromR2("categories", category.PhotoURL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete category image: %v", err)
		}
	}

	// Set PhotoURL ke kosong
	category.PhotoURL = ""
	if err := tx.Save(&category).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetCategoryOptions() ([]dto.CategoryOption, error) {
	var categories []models.Category
	if err := config.DB.Select("uuid, name").
		Where("is_published = ?", true).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	options := make([]dto.CategoryOption, len(categories))
	for i, category := range categories {
		options[i] = dto.CategoryOption{
			UUID: category.UUID,
			Name: category.Name,
		}
	}

	return options, nil
}

func GetCategoryBySlug(slug string) (dto.CategoryBySlugResponse, error) {
	var category models.Category
	if err := config.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		return dto.CategoryBySlugResponse{}, err
	}

	if category.UUID == "" {
		return dto.CategoryBySlugResponse{}, fmt.Errorf("category not found")
	}

	var users []struct {
		UUID string
		Slug string
		Name string
	}

	subQuery := config.DB.
		Table("albums").
		Select("user_id").
		Where("category_id = ? AND is_published = ?", category.ID, true)

	if err := config.DB.
		Table("users").
		Select("uuid, slug, name").
		Where("id IN (?) AND is_published = ?", subQuery, true).
		Scan(&users).Error; err != nil {
		return dto.CategoryBySlugResponse{}, err
	}

	// You can attach users to category or return them as needed
	// Example: category.Users = users (if you have such a field)

	usersResp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		usersResp = append(usersResp, dto.UserResponse{
			UUID: u.UUID,
			Name: u.Name,
			Slug: u.Slug,
		})
	}

	return dto.CategoryBySlugResponse{
		UUID:        category.UUID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		PhotoUrl:    category.PhotoURL,
		Users:       usersResp,
	}, nil
}
