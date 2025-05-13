package controllers

import (
	"net/http"
	"strconv"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func GetFaqs(c *gin.Context) {
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

	faqs, total, err := services.GetAllFaqs(pageInt, limitInt, search)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get faqs")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data":  faqs,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})

}

func CreateFaq(c *gin.Context) {
	var input services.FaqInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := validate.Struct(&input); err != nil {
		// Bisa custom format error jika mau
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	faq, err := services.CreateFaq(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": faq})
}

func EditFaq(c *gin.Context) {
	id := c.Param("uuid")

	_, err := services.GetFaqByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get faq")
		return
	}

	var input services.FaqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := validate.Struct(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedFaq, err := services.UpdateFaq(id, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": updatedFaq})
}

func DeleteFaq(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	err := services.DeleteFaq(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "deleted successfully",
	})
}

func GetFaqByUUID(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	faq, err := services.GetFaqByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if faq.UUID == "" {
		utils.RespondError(c, http.StatusNotFound, "faq not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": gin.H{
			"uuid":         faq.UUID,
			"answer_en":    faq.AnswerEn,
			"answer_id":    faq.AnswerID,
			"question_en":  faq.QuestionEn,
			"question_id":  faq.QuestionID,
			"is_published": faq.IsPublished,
			"created_at":   faq.CreatedAt,
			"updated_at":   faq.UpdatedAt,
		},
	})
}
