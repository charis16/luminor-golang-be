package controllers

import (
	"net/http"
	"strconv"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetPublishedCategories(c *gin.Context) {
	categories, err := services.GetPublishedCategories()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get faqs")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": categories,
	})
}

func GetCategories(c *gin.Context) {
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

	categories, total, err := services.GetAllCategories(pageInt, limitInt, search)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get categories")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data":  categories,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})

}

func CreateCategory(c *gin.Context) {
	var input services.CategoryInput

	// Gunakan ShouldBind
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if err := validate.Struct(&input); err != nil {
		// Bisa custom format error jika mau
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if fileHeader, err := c.FormFile("image"); err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to open thumbnail")
			return
		}
		defer file.Close()

		url, err := utils.UploadToR2(file, fileHeader, "categories") // thumbnail juga simpan ke prefix albums
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to upload thumbnail")
			return
		}
		input.PhotoUrl = url
	}

	category, err := services.CreateCategory(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": category})
}

func EditCategory(c *gin.Context) {
	id := c.Param("uuid")

	_, err := services.GetCategoryByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get category")
		return
	}

	var input services.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := validate.Struct(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedCategory, err := services.UpdateCategory(id, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": updatedCategory})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	err := services.DeleteCategory(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "user deleted successfully",
	})
}

func GetCategoryByUUID(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	category, err := services.GetCategoryByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if category.UUID == "" {
		utils.RespondError(c, http.StatusNotFound, "category not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": gin.H{
			"uuid":         category.UUID,
			"name":         category.Name,
			"description":  category.Description,
			"slug":         category.Slug,
			"photo_url":    category.PhotoURL,
			"is_published": category.IsPublished,
			"created_at":   category.CreatedAt,
			"updated_at":   category.UpdatedAt,
		},
	})
}

func DeleteImageCategory(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	err := services.DeleteImageCategory(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "image user deleted successfully",
	})
}

func GetCategoryOptions(c *gin.Context) {
	options, err := services.GetCategoryOptions()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get category options")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": options,
	})
}

func GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if slug == "" {
		utils.RespondError(c, http.StatusBadRequest, "slug is required")
		return
	}

	category, err := services.GetCategoryBySlug(slug)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if category.UUID == "" {
		utils.RespondError(c, http.StatusNotFound, "category not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": gin.H{
			"uuid":        category.UUID,
			"name":        category.Name,
			"description": category.Description,
			"slug":        category.Slug,
			"photo_url":   category.PhotoUrl,
			"users":       category.Users,
		},
	})
}
