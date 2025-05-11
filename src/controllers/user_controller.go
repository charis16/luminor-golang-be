package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	urlTikTok := c.PostForm("url_tikTok")
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

	input := services.UserInput{
		Name:         name,
		Email:        email,
		Role:         role,
		Description:  description,
		PhotoURL:     photoURL,
		Password:     password,
		URLInstagram: urlInstagram,
		URLTikTok:    urlTikTok,
		URLFacebook:  urlFacebook,
		URLYoutube:   urlYoutube,
		PhoneNumber:  phoneNumber,
		IsPublished:  isPublished,
	}

	user, err := services.CreateUser(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": user,
	})
}

func EditUser(c *gin.Context) {
	id := c.Param("uuid")
	name := c.PostForm("name")
	email := c.PostForm("email")
	role := c.PostForm("role")
	description := c.PostForm("description")
	password := c.PostForm("password")
	urlInstagram := c.PostForm("url_instagram")
	urlTikTok := c.PostForm("url_tikTok")
	urlFacebook := c.PostForm("url_facebook")
	urlYoutube := c.PostForm("url_youtube")
	phoneNumber := c.PostForm("phone_number")
	isPublished := c.PostForm("is_published")

	user, err := services.GetUserByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	if name == "" || email == "" || role == "" {
		utils.RespondError(c, http.StatusBadRequest, "name, email, and role are required")
		return
	}

	var photoURL string

	fileHeader, err := c.FormFile("photo")
	if err == nil {
		// File is provided, upload the new photo
		file, err := fileHeader.Open()
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to open uploaded file")
			return
		}
		defer file.Close()

		photoURL, err = utils.UploadToMinio("users", file, fileHeader)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to upload photo")
			return
		}

		if user.Photo != "" {
			err = utils.DeleteFromMinio("users", user.Photo)
			if err != nil {
				fmt.Println("Warning: failed to delete old photo:", err)
			}
		}
	} else {
		// No new file provided, use the existing photo URL
		photoURL = user.Photo
	}

	input := services.UserInput{
		Name:         name,
		Email:        email,
		Role:         role,
		Description:  description,
		PhotoURL:     photoURL,
		Password:     password,
		URLInstagram: urlInstagram,
		URLTikTok:    urlTikTok,
		URLFacebook:  urlFacebook,
		URLYoutube:   urlYoutube,
		PhoneNumber:  phoneNumber,
		IsPublished:  isPublished,
	}

	user, err = services.UpdateUser(id, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	err := services.DeleteUser(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "user deleted successfully",
	})
}

func GetUserByUUID(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	user, err := services.GetUserByUUID(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if user.UUID == "" {
		utils.RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": gin.H{
			"uuid":          user.UUID,
			"name":          user.Name,
			"email":         user.Email,
			"role":          user.Role,
			"description":   user.Description,
			"phone_number":  user.PhoneNumber,
			"url_instagram": user.URLInstagram,
			"url_tiktok":    user.URLTiktok,
			"url_facebook":  user.URLFacebook,
			"url_youtube":   user.URLYoutube,
			"is_published":  user.IsPublished,
			"photo_url":     user.Photo, // final accessible URL
		},
	})
}

func ProxyUserImage(c *gin.Context) {
	filename := c.Query("filename")
	utils.StreamImageFromMinio(c, "users", filename, "image/*", 5*time.Minute)
}
