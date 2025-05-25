package controllers

import (
	"net/http"
	"strconv"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetAlbums(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	albums, total, err := services.GetAllAlbums(pageInt, limitInt, search)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get albums")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data":  albums,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})

}

func CreateAlbum(c *gin.Context) {
	var input services.AlbumInput

	// Gunakan ShouldBind
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	// Upload images
	form, err := c.MultipartForm()
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid multipart form")
		return
	}

	files := form.File["images"]
	var imageUrls []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to open file")
			return
		}
		defer file.Close()

		url, err := utils.UploadToMinio("albums", file, fileHeader)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to upload")
			return
		}

		imageUrls = append(imageUrls, url)
	}
	input.Images = imageUrls

	// Upload thumbnail (opsional)
	if fileHeader, err := c.FormFile("thumbnail"); err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to open thumbnail")
			return
		}
		defer file.Close()

		url, err := utils.UploadToMinio("albums", file, fileHeader)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to upload thumbnail")
			return
		}
		input.Thumbnail = url
	}

	// Validasi pakai validator
	if err := validate.Struct(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	album, err := services.CreateAlbum(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": album})
}

func EditAlbum(c *gin.Context) {
	id := c.Param("uuid")
	// Cek apakah album dengan UUID tersebut ada
	_, err := services.GetAlbumByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to get album")
		return
	}

	var input services.AlbumInput

	// Ambil field dari PostForm
	input.Slug = c.PostForm("slug")
	input.Title = c.PostForm("title")
	input.CategoryId = c.PostForm("category_id")
	input.Description = c.PostForm("description")
	input.UserID = c.PostForm("user_id")
	input.IsPublished = c.PostForm("is_published")

	// Ambil file-file jika ada
	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		utils.RespondError(c, http.StatusBadRequest, "Invalid form data")
		return
	}

	var imageUrls []string
	if form != nil && form.File != nil {
		files := form.File["images"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				utils.RespondError(c, http.StatusInternalServerError, "Failed to open uploaded file")
				return
			}
			defer file.Close()

			url, err := utils.UploadToMinio("albums", file, fileHeader)
			if err != nil {
				utils.RespondError(c, http.StatusInternalServerError, "Failed to upload image")
				return
			}

			imageUrls = append(imageUrls, url)
		}
	}
	input.Images = imageUrls

	if fileHeader, err := c.FormFile("thumbnail"); err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to open thumbnail")
			return
		}
		defer file.Close()

		url, err := utils.UploadToMinio("albums", file, fileHeader)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to upload thumbnail")
			return
		}
		input.Thumbnail = url
	}

	// Validasi input
	if err := validate.Struct(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update ke service
	updatedAlbum, err := services.UpdateAlbum(id, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": updatedAlbum})
}

func DeleteAlbum(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	err := services.DeleteAlbum(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "deleted successfully",
	})
}

func GetAlbumByUUID(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	album, err := services.GetAlbumByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if album.UUID == "" {
		utils.RespondError(c, http.StatusNotFound, "faq not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": gin.H{
			"uuid":         album.UUID,
			"slug":         album.Slug,
			"title":        album.Title,
			"category_id":  album.Category.UUID,
			"description":  album.Description,
			"thumbnail":    album.Thumbnail,
			"images":       album.Images,
			"user_id":      album.User.UUID,
			"is_published": album.IsPublished,
			"created_at":   album.CreatedAt,
			"updated_at":   album.UpdatedAt,
		},
	})
}
