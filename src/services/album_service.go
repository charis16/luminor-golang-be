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

type AlbumInput struct {
	Slug        string   `form:"slug" binding:"required"`
	Title       string   `form:"title" binding:"required"`
	CategoryId  string   `form:"category_id" binding:"required"`
	Description string   `form:"description" binding:"required"`
	UserID      string   `form:"user_id" binding:"required"`
	IsPublished string   `form:"is_published" binding:"required"`
	Images      []string `form:"-"` // handled manually
	Thumbnail   string   `form:"-"` // handled manually
}

func GetAllAlbums(page int, limit int, search string) ([]dto.AlbumResponse, int64, error) {
	var albums []models.Album
	var total int64

	query := config.DB.Model(&models.Album{})

	// Apply search filter if search term is provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Select("id", "uuid", "slug", "title", "category_id", "description", "images", "thumbnail", "is_published", "user_id", "created_at", "updated_at").
		Preload("User").
		Preload("Category").
		Limit(limit).
		Offset(offset).
		Find(&albums).Error; err != nil {
		return nil, 0, err
	}

	// Mapping ke response DTO
	response := make([]dto.AlbumResponse, len(albums))
	for i, album := range albums {

		response[i] = dto.AlbumResponse{
			UUID:         album.UUID,
			Slug:         album.Slug,
			Title:        album.Title,
			CategoryId:   album.Category.UUID,
			CategoryName: album.Category.Name,
			Description:  album.Description,
			Images:       album.Images,
			Thumbnail:    album.Thumbnail,
			IsPublished:  album.IsPublished,
			CreatedAt:    album.CreatedAt,
			UpdatedAt:    album.UpdatedAt,
			UserID:       album.User.UUID,
			UserName:     album.User.Name,
			UserAvatar:   album.User.Photo,
		}
	}

	return response, total, nil
}

func CreateAlbum(input AlbumInput) (*models.Album, error) {
	tx := config.DB.Begin() // Mulai transaksi

	category, err := GetCategoryByUUID(input.CategoryId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user, err := GetUserByUUID(input.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Cek apakah slug sudah ada
	var existingAlbum models.Album
	if err := tx.Where("slug = ?", input.Slug).First(&existingAlbum).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("slug already exists")
	}
	album := models.Album{
		Slug:        input.Slug,
		Title:       input.Title,
		CategoryID:  category.ID,
		Description: input.Description,
		Images:      input.Images,
		Thumbnail:   input.Thumbnail,
		UserID:      user.ID,
		IsPublished: input.IsPublished == "true",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&album).Error; err != nil {
		tx.Rollback() // rollback jika error
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err // commit gagal
	}

	return &album, nil
}

func GetAlbumByUUID(uuid string) (models.Album, error) {
	var album models.Album
	if err := config.DB.
		Select("id", "uuid", "slug", "title", "category_id", "description", "images", "thumbnail", "is_published", "user_id", "created_at", "updated_at").
		Preload("User").
		Preload("Category").
		Where("uuid = ?", uuid).First(&album).Error; err != nil {
		return models.Album{}, err
	}
	return album, nil
}

func UpdateAlbum(uuid string, input AlbumInput) (models.Album, error) {
	tx := config.DB.Begin()
	var album models.Album

	album, err := GetAlbumByUUID(uuid)
	if err != nil {
		tx.Rollback()
		return models.Album{}, err
	}

	if input.Slug != album.Slug {
		var existingAlbum models.Album
		if err := tx.Where("slug = ? AND uuid != ?", input.Slug, uuid).First(&existingAlbum).Error; err == nil {
			tx.Rollback()
			return models.Album{}, fmt.Errorf("slug already exists")
		}
	}

	category, err := GetCategoryByUUID(input.CategoryId)
	if err != nil {
		tx.Rollback()
		return models.Album{}, err
	}

	user, err := GetUserByUUID(input.UserID)
	if err != nil {
		tx.Rollback()
		return models.Album{}, err
	}

	album.Slug = input.Slug
	album.Title = input.Title
	album.Description = input.Description

	// Update Category if changed
	if input.CategoryId != "" && category.ID != album.CategoryID {
		album.CategoryID = category.ID
	}

	if input.Images != nil {
		var filteredOldImages []string
		for _, img := range album.Images {
			img = strings.Trim(img, `"`)
			if img != "" {
				filteredOldImages = append(filteredOldImages, img)
			}
		}
		combinedImages := append(filteredOldImages, input.Images...)
		album.Images = combinedImages
	}

	if input.Thumbnail != "" {
		album.Thumbnail = input.Thumbnail
	}

	album.UserID = user.ID
	album.IsPublished = input.IsPublished == "true"
	album.UpdatedAt = time.Now()

	if err := tx.Save(&album).Error; err != nil {
		tx.Rollback()
		return models.Album{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.Album{}, err
	}

	return album, nil
}

func DeleteAlbum(uuid string) error {
	tx := config.DB.Begin()

	// Ambil album untuk mendapatkan daftar images dan thumbnail
	album, err := GetAlbumByUUID(uuid)
	if err != nil {
		return err
	}

	// Hapus images dari MinIO
	for _, img := range album.Images {
		img = strings.Trim(img, `"`)
		if img != "" {
			if err := utils.DeleteFromMinio("albums", img); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to delete album image: %v", err)
			}
		}
	}

	// Hapus thumbnail dari MinIO
	if album.Thumbnail != "" {
		if err := utils.DeleteFromMinio("albums", album.Thumbnail); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete album image: %v", err)
		}
	}

	// Hapus album dari database
	if err := config.DB.Where("uuid = ?", uuid).Delete(&models.Album{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteImageFromAlbum(uuid string, imageURL string) error {
	tx := config.DB.Begin()

	// Ambil album
	album, err := GetAlbumByUUID(uuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Trim input
	imageURL = strings.Trim(imageURL, `"`)

	// Ambil nama file dari URL untuk delete dari Minio
	parts := strings.Split(imageURL, "/")
	imageFilename := parts[len(parts)-1]

	// Delete from Minio
	if imageFilename != "" {
		if err := utils.DeleteFromMinio("albums", imageFilename); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete album image from MinIO: %v", err)
		}
	}

	// Filter out URL yang dihapus
	var updatedImages []string
	for _, img := range album.Images {
		if img != imageURL { // URL match
			updatedImages = append(updatedImages, img)
		}
	}

	// Update album
	album.Images = updatedImages
	if err := tx.Save(&album).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
