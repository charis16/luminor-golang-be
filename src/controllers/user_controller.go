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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	users, total, err := services.GetAllUsers(pageInt, limitInt, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

func CreateUser(c *gin.Context) {
	fmt.Println("Content-Type:", c.Request.Header.Get("Content-Type"))

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, and role are required"})
		return
	}

	var photoURL string

	fileHeader, err := c.FormFile("photo")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
			return
		}
		defer file.Close()

		fmt.Println("Uploading photo to Minio...")
		photoURL, err = utils.UploadToMinio("users", file, fileHeader)
		if err != nil {
			fmt.Println("Failed to upload photo:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload photo"})
			return
		}
		fmt.Println("Photo uploaded successfully. URL:", photoURL)
	} else {
		fmt.Println("No photo uploaded.")
	}

	user, err := services.CreateUser(name, email, role, description, photoURL, password, urlInstagram, url_tikTok, urlFacebook, urlYoutube, phoneNumber, isPublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}
