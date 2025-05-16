package controllers

import (
	"net/http"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func GetWebsite(c *gin.Context) {

	website, _, err := services.GetWebsite()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get website")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": website,
	})

}

func CreateWebsiteInformation(c *gin.Context) {
	var input services.WebsiteInput

	if err := c.ShouldBind(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := validate.Struct(&input); err != nil {
		// Bisa custom format error jika mau
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	faq, err := services.CreateWebsiteInformation(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": faq})
}

func EditWebsiteInformation(c *gin.Context) {
	id := c.Param("uuid")

	_, err := services.GetFaqByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get faq")
		return
	}

	var input services.WebsiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := validate.Struct(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedFaq, err := services.EditWebsiteInformation(id, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": updatedFaq})
}
