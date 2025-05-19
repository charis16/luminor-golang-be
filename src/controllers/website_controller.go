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
		utils.RespondSuccess(c, gin.H{
			"data": nil,
		})
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": website,
	})

}
func CreateWebsiteInformation(c *gin.Context) {
	var input services.WebsiteInput

	contentType := c.GetHeader("Content-Type")
	if contentType != "" && (contentType == "multipart/form-data" || len(contentType) > 19 && contentType[:19] == "multipart/form-data") {
		metaTitle := c.PostForm("meta_title")
		metaKeywords := c.PostForm("meta_keywords")
		metaDescription := c.PostForm("meta_description")

		var ogImage string

		fileHeader, err := c.FormFile("ogImage")
		if err == nil && fileHeader != nil {
			file, openErr := fileHeader.Open()
			if openErr != nil {
				utils.RespondError(c, http.StatusInternalServerError, "failed to open uploaded file")
				return
			}
			defer file.Close()

			ogImage, err = utils.UploadToMinio("websites", file, fileHeader)
			if err != nil {
				utils.RespondError(c, http.StatusInternalServerError, "failed to upload photo")
				return
			}
		}

		input := services.WebsiteInput{}

		if metaTitle != "" {
			input.MetaTitle = metaTitle
		}
		if metaKeywords != "" {
			input.MetaKeyword = metaKeywords
		}
		if metaDescription != "" {
			input.MetaDesc = metaDescription
		}
		if ogImage != "" {
			input.OgImage = ogImage
		}

	} else {
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.RespondError(c, http.StatusBadRequest, "Invalid input format: "+err.Error())
			return
		}
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
	var input services.WebsiteInput

	_, err := services.GetWebsiteByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get website information")
		return
	}

	contentType := c.GetHeader("Content-Type")
	if contentType != "" && (contentType == "multipart/form-data" || len(contentType) > 19 && contentType[:19] == "multipart/form-data") {
		metaTitle := c.PostForm("meta_title")
		metaKeywords := c.PostForm("meta_keywords")
		metaDescription := c.PostForm("meta_description")

		var ogImage string

		fileHeader, err := c.FormFile("ogImage")
		if err == nil && fileHeader != nil {
			file, openErr := fileHeader.Open()
			if openErr != nil {
				utils.RespondError(c, http.StatusInternalServerError, "failed to open uploaded file")
				return
			}
			defer file.Close()

			ogImage, err = utils.UploadToMinio("websites", file, fileHeader)
			if err != nil {
				utils.RespondError(c, http.StatusInternalServerError, "failed to upload photo")
				return
			}
		}

		input := services.WebsiteInput{}

		if metaTitle != "" {
			input.MetaTitle = metaTitle
		}
		if metaKeywords != "" {
			input.MetaKeyword = metaKeywords
		}
		if metaDescription != "" {
			input.MetaDesc = metaDescription
		}
		if ogImage != "" {
			input.OgImage = ogImage
		}

	} else {
		if err := c.ShouldBindJSON(&input); err != nil {
			utils.RespondError(c, http.StatusBadRequest, "Invalid input format: "+err.Error())
			return
		}
	}

	updatedFaq, err := services.EditWebsiteInformation(id, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": updatedFaq})
}
