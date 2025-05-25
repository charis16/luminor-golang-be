package controllers

import (
	"net/http"
	"strconv"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

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

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := validate.Struct(&input); err != nil {
		// Bisa custom format error jika mau
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
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
			"is_published": category.IsPublished,
			"created_at":   category.CreatedAt,
			"updated_at":   category.UpdatedAt,
		},
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
