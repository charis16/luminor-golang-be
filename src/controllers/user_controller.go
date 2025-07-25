package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func GetUserPortfolioBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.RespondError(c, http.StatusBadRequest, "slug is required")
		return
	}

	user, err := services.GetUserPortfolioBySlug(slug)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to get user portfolio by slug")
		return
	}

	if user.User.UUID == "" {
		utils.RespondError(c, http.StatusNotFound, "album not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": gin.H{
			"user": gin.H{
				"uuid":          user.User.UUID,
				"name":          user.User.Name,
				"email":         user.User.Email,
				"role":          user.User.Role,
				"description":   user.User.Description,
				"slug":          user.User.Slug,
				"phone_number":  user.User.PhoneNumber,
				"url_instagram": user.User.URLInstagram,
				"url_tiktok":    user.User.URLTikTok,
				"url_facebook":  user.User.URLFacebook,
				"url_youtube":   user.User.URLYoutube,
				"is_published":  user.User.IsPublished,
				"photo_url": func() interface{} {
					if user.User.Photo == "" {
						return nil
					} else {
						return user.User.Photo
					}
				}(),
			},
			"categories": user.Categories,
		},
	})
}

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

	if name == "" || email == "" || role == "" {
		utils.RespondError(c, http.StatusBadRequest, "name, email, and role are required")
		return
	}

	slug := c.PostForm("slug")
	description := c.PostForm("description")
	password := c.PostForm("password")
	urlInstagram := c.PostForm("url_instagram")
	urlTikTok := c.PostForm("url_tikTok")
	urlFacebook := c.PostForm("url_facebook")
	urlYoutube := c.PostForm("url_youtube")
	phoneNumber := c.PostForm("phone_number")
	isPublished := c.PostForm("is_published")
	canLogin := c.PostForm("can_login")

	var photoURL string

	fileHeader, err := c.FormFile("photo")
	if err == nil && fileHeader != nil {
		file, openErr := fileHeader.Open()
		if openErr != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to open uploaded file")
			return
		}
		defer file.Close()

		photoURL, err = utils.UploadToR2(file, fileHeader, "users")
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to upload photo")
			return
		}
	}

	input := services.UserInput{
		Slug:         slug,
		Name:         name,
		Email:        email,
		Role:         role,
		Description:  description,
		Password:     password,
		URLInstagram: urlInstagram,
		URLTikTok:    urlTikTok,
		URLFacebook:  urlFacebook,
		URLYoutube:   urlYoutube,
		PhoneNumber:  phoneNumber,
		IsPublished:  isPublished == "true",
		CanLogin:     canLogin == "true",
		PhotoURL:     photoURL,
	}

	user, err := services.CreateUser(input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"data": user})
}

func EditUser(c *gin.Context) {
	id := c.Param("uuid")
	slug := c.PostForm("slug")
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
	canLogin := c.PostForm("can_login")

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

		photoURL, err = utils.UploadToR2(file, fileHeader, "users")
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "failed to upload photo")
			return
		}

		if user.Photo != "" {
			err = utils.DeleteFromR2("users", user.Photo)
			if err != nil {
				fmt.Println("Warning: failed to delete old photo:", err)
			}
		}
	} else {
		// No new file provided, use the existing photo URL
		photoURL = user.Photo
	}

	input := services.UserInput{
		Slug:         slug,
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
		IsPublished:  isPublished == "true",
		CanLogin:     canLogin == "true",
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
			"slug":          user.Slug,
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
			"can_login": func() interface{} {
				if user.Password == "" {
					return false
				} else {
					return true
				}
			}(), // final accessible URL
			"photo_url": func() interface{} {
				if user.Photo == "" {
					return nil
				} else {
					return user.Photo
				}
			}(), // final accessible URL
		},
	})
}

func DeleteImageUser(c *gin.Context) {
	id := c.Param("uuid")

	if id == "" {
		utils.RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	err := services.DeleteImageUser(id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "image user deleted successfully",
	})
}

func GetUserOptions(c *gin.Context) {
	options, err := services.GetUserOptions()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get user options")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": options,
	})
}

func GetTeamMembers(c *gin.Context) {
	teamMembers, err := services.GetTeamMembers()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get team members")
		return
	}

	if len(teamMembers) == 0 {
		utils.RespondSuccess(c, gin.H{
			"data": []string{},
		})
		return
	}

	utils.RespondSuccess(c, gin.H{
		"data": func() []gin.H {
			result := make([]gin.H, 0, len(teamMembers))
			for _, user := range teamMembers {
				result = append(result, gin.H{
					"uuid":          user.UUID,
					"name":          user.Name,
					"email":         user.Email,
					"role":          user.Role,
					"description":   user.Description,
					"slug":          user.Slug,
					"phone_number":  user.PhoneNumber,
					"url_instagram": user.URLInstagram,
					"url_tiktok":    user.URLTikTok,
					"url_facebook":  user.URLFacebook,
					"url_youtube":   user.URLYoutube,
					"is_published":  user.IsPublished,
					"photo_url": func() interface{} {
						if user.Photo == "" {
							return nil
						} else {
							return user.Photo
						}
					}(),
				})
			}
			return result
		}(),
	})
}
