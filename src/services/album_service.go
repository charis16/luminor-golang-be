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
	Slug        string   `form:"slug"`
	Title       string   `form:"title" binding:"required"`
	CategoryId  string   `form:"category_id" binding:"required"`
	Description string   `form:"description" binding:"required"`
	UserID      string   `form:"user_id" binding:"required"`
	IsPublished string   `form:"is_published" binding:"required"`
	YoutubeURL  string   `form:"youtube_url"`
	Images      []string `form:"-"` // handled manually
	Thumbnail   string   `form:"-"` // handled manually
}

type DeleteImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
}

func mapAlbumToDTO(album models.Album) dto.AlbumResponse {
	return dto.AlbumResponse{
		UUID:         album.UUID,
		Slug:         album.Slug,
		Title:        album.Title,
		CategoryId:   album.Category.UUID,
		CategoryName: album.Category.Name,
		CategorySlug: album.Category.Slug,
		Description:  album.Description,
		YoutubeURL:   album.YoutubeURL,
		Images:       album.Images,
		Thumbnail:    album.Thumbnail,
		IsPublished:  album.IsPublished,
		CreatedAt:    album.CreatedAt,
		UpdatedAt:    album.UpdatedAt,
		UserID:       album.User.UUID,
		UserName:     album.User.Name,
		UserAvatar:   album.User.Photo,
		UserSlug:     album.User.Slug,
	}
}

func GetLatestAlbums() ([]dto.AlbumResponse, error) {
	var albums []models.Album
	if err := config.DB.
		Preload("User").
		Preload("Category").
		Order("created_at DESC").
		Where("is_published = ?", true).
		Limit(20).
		Find(&albums).Error; err != nil {
		return []dto.AlbumResponse{}, err
	}

	var items []dto.AlbumResponse
	for _, album := range albums {
		items = append(items, mapAlbumToDTO(album))
	}

	return items, nil
}

func GetAlbumByCategorySlug(slug string, nextTime int,
	limit int, filter string) (dto.AlbumResponseList, error) {
	var albums []models.Album
	var nextTimestamp int64

	query := config.DB.
		Preload("User").
		Preload("Category").
		Where("is_published = ?", true).
		Order("created_at DESC").
		Limit(limit)

	if slug != "" && slug != "all" {
		var category models.Category
		if err := config.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
			return dto.AlbumResponseList{
				Data:      []dto.AlbumResponse{},
				NextValue: 999999999,
			}, err
		}
		query = query.Where("category_id = ?", category.ID)
	}

	if filter != "" && filter != "all" {
		var user models.User
		if err := config.DB.Where("slug = ?", filter).First(&user).Error; err != nil {
			return dto.AlbumResponseList{
				Data:      []dto.AlbumResponse{},
				NextValue: 999999999,
			}, err
		}
		query = query.Where("user_id = ?", user.ID)
	}

	if nextTime != 0 {
		query = query.Where("UNIX_TIMESTAMP(created_at) < ?", nextTime)
	}

	if err := query.Find(&albums).Error; err != nil {
		return dto.AlbumResponseList{
			Data:      []dto.AlbumResponse{},
			NextValue: 999999999,
		}, err
	}

	if len(albums) == 0 {
		return dto.AlbumResponseList{
			Data:      []dto.AlbumResponse{},
			NextValue: 999999999,
		}, nil
	}

	nextTimestamp = albums[len(albums)-1].CreatedAt.Unix()

	var items []dto.AlbumResponse
	for _, album := range albums {
		items = append(items, mapAlbumToDTO(album))
	}

	return dto.AlbumResponseList{
		Data:      items,
		NextValue: nextTimestamp,
	}, nil
}

func GetDetailAlbumBySlug(slug string) (dto.AlbumResponse, error) {
	var album models.Album
	if err := config.DB.
		Preload("User").
		Preload("Category").
		Where("slug = ? AND is_published = ?", slug, true).
		First(&album).Error; err != nil {
		return dto.AlbumResponse{}, err
	}

	return mapAlbumToDTO(album), nil
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
		response[i] = mapAlbumToDTO(album)
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

	slug := input.Slug
	if slug == "" {
		slug = utils.GenerateSlug(input.Title)
	} else {
		slug = utils.GenerateSlug(slug)
	}

	// Cek apakah slug sudah ada
	var existingAlbum models.Album
	if err := tx.Where("slug = ?", slug).First(&existingAlbum).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("slug already exists")
	}
	album := models.Album{
		Slug:        slug,
		Title:       input.Title,
		CategoryID:  category.ID,
		Description: input.Description,
		Images:      input.Images,
		Thumbnail:   input.Thumbnail,
		YoutubeURL:  input.YoutubeURL,
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

	slug := utils.GenerateSlug(input.Slug)
	if slug != album.Slug {
		var existingAlbum models.Album
		if err := tx.Where("slug = ? AND uuid != ?", slug, uuid).First(&existingAlbum).Error; err == nil {
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

	album.Slug = slug
	album.Title = input.Title
	album.Description = input.Description
	album.YoutubeURL = input.YoutubeURL

	// Update Category if changed
	if input.CategoryId != "" && category.ID != album.CategoryID {
		album.CategoryID = category.ID
	}

	if input.Images != nil {
		combined := append(album.Images, input.Images...)
		album.Images = utils.RemoveDuplicateStrings(combined)
	}

	if input.Thumbnail != "" && input.Thumbnail != "undefined" {
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
			if err := utils.DeleteFromR2("albums", img); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to delete album image: %v", err)
			}
		}
	}

	// Hapus thumbnail dari MinIO
	if album.Thumbnail != "" {
		if err := utils.DeleteFromR2("albums", album.Thumbnail); err != nil {
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

	// Ambil nama file dari URL
	parts := strings.Split(imageURL, "/")
	imageFilename := parts[len(parts)-1]

	// Hapus dari MinIO
	if imageFilename != "" {
		if err := utils.DeleteFromR2("albums", imageFilename); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete album image from MinIO: %v", err)
		}
	}

	// Cek apakah imageURL adalah thumbnail
	if album.Thumbnail == imageURL {
		album.Thumbnail = ""
	} else {
		// Kalau bukan thumbnail, hapus dari album.Images
		var updatedImages []string
		for _, img := range album.Images {
			if img != imageURL {
				updatedImages = append(updatedImages, img)
			}
		}
		album.Images = updatedImages
	}

	// Simpan perubahan
	if err := tx.Save(&album).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
