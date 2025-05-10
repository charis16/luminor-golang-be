package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
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

	users, total, err := services.GetAllUsers(pageInt, limitInt, search)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get users")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data":  users,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})

}

func CreateUser(c *gin.Context) {

	name := c.PostForm("name")
	email := c.PostForm("email")
	role := c.PostForm("role")
	description := c.PostForm("description")
	password := c.PostForm("password")
	urlInstagram := c.PostForm("url_instagram")
	url_tikTok := c.PostForm("url_tikTok")
	urlFacebook := c.PostForm("url_facebook")
	urlYoutube := c.PostForm("url_youtube")
	phoneNumber := c.PostForm("phone_number")
	isPublished := c.PostForm("is_published")

	if name == "" || email == "" || role == "" {
		utils.RespondError(c, http.StatusBadRequest, "name, email, and role are required")
		return
	}

	var photoURL string

	fileHeader, err := c.FormFile("photo")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to open uploaded file")
			return
		}
		defer file.Close()

		fmt.Println("Uploading photo to Minio...")
		photoURL, err = utils.UploadToMinio("users", file, fileHeader)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to upload photo")
			return
		}
	}

	user, err := services.CreateUser(name, email, role, description, photoURL, password, urlInstagram, url_tikTok, urlFacebook, urlYoutube, phoneNumber, isPublished)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": user,
	})
}
