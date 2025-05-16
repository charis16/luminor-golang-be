package services

import (
	"time"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
)

type WebsiteInput struct {
	Address            string `json:"address"`
	PhoneNumber        string `json:"phone_number"`
	Latitude           string `json:"latitude"`
	Longitude          string `json:"longitude"`
	Email              string `json:"email"`
	UrlInstagram       string `json:"url_instagram"`
	UrlTikTok          string `json:"url_tiktok"`
	AboutUsBriefHomeEn string `json:"about_us_brief_home_en"`
	AboutUsEn          string `json:"about_us_en"`
	AboutUsID          string `json:"about_us_id"`
	AboutUsBriefHomeID string `json:"about_us_brief_home_id"`
	VideoWeb           string `json:"video_web"`
	VideoMobile        string `json:"video_mobile"`
	MetaTitle          string `json:"meta_title"`
	MetaDesc           string `json:"meta_desc"`
	MetaKeyword        string `json:"meta_keyword"`
	OgImage            string `json:"og_image"`
}

func GetWebsite() ([]dto.WebsiteResponse, int64, error) {
	var website models.Website
	if err := config.DB.First(&website).Error; err != nil {
		return nil, 0, err
	}

	response := []dto.WebsiteResponse{
		{
			UUID:               website.UUID,
			AboutUsBriefHomeEn: website.AboutUsBriefHomeEn,
			AboutUsEn:          website.AboutUsEn,
			AboutUsID:          website.AboutUsID,
			AboutUsBriefHomeID: website.AboutUsBriefHomeID,
			VideoWeb:           website.VideoWeb,
			VideoMobile:        website.VideoMobile,
			MetaTitle:          website.MetaTitle,
			MetaDesc:           website.MetaDesc,
			MetaKeyword:        website.MetaKeyword,
			OgImage:            website.OgImage,
			Address:            website.Address,
			PhoneNumber:        website.PhoneNumber,
			Latitude:           website.Latitude,
			Longitude:          website.Longitude,
			Email:              website.Email,
			UrlInstagram:       website.URLInstagram,
			UrlTikTok:          website.URLFacebook,
			CreatedAt:          website.CreatedAt,
			UpdatedAt:          website.UpdatedAt,
		},
	}

	return response, 1, nil
}

func CreateWebsiteInformation(input WebsiteInput) (*models.Website, error) {
	tx := config.DB.Begin() // Mulai transaksi
	website := models.Website{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if input.Address != "" {
		website.Address = input.Address
	}
	if input.PhoneNumber != "" {
		website.PhoneNumber = input.PhoneNumber
	}
	if input.Latitude != "" {
		website.Latitude = input.Latitude
	}
	if input.Longitude != "" {
		website.Longitude = input.Longitude
	}
	if input.Email != "" {
		website.Email = input.Email
	}
	if input.UrlInstagram != "" {
		website.URLInstagram = input.UrlInstagram
	}
	if input.UrlTikTok != "" {
		website.URLFacebook = input.UrlTikTok
	}
	if input.AboutUsBriefHomeEn != "" {
		website.AboutUsBriefHomeEn = input.AboutUsBriefHomeEn
	}
	if input.AboutUsEn != "" {
		website.AboutUsEn = input.AboutUsEn
	}
	if input.AboutUsID != "" {
		website.AboutUsID = input.AboutUsID
	}
	if input.AboutUsBriefHomeID != "" {
		website.AboutUsBriefHomeID = input.AboutUsBriefHomeID
	}
	if input.VideoWeb != "" {
		website.VideoWeb = input.VideoWeb
	}
	if input.VideoMobile != "" {
		website.VideoMobile = input.VideoMobile
	}
	if input.MetaTitle != "" {
		website.MetaTitle = input.MetaTitle
	}
	if input.MetaDesc != "" {
		website.MetaDesc = input.MetaDesc
	}
	if input.MetaKeyword != "" {
		website.MetaKeyword = input.MetaKeyword
	}
	if input.OgImage != "" {
		website.OgImage = input.OgImage
	}

	if err := tx.Create(&website).Error; err != nil {
		tx.Rollback() // rollback jika error
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err // commit gagal
	}

	return &website, nil
}

func EditWebsiteInformation(uuid string, input WebsiteInput) (models.Website, error) {
	tx := config.DB.Begin()
	var website models.Website

	if err := tx.Where("uuid = ?", uuid).First(&website).Error; err != nil {
		tx.Rollback()
		return models.Website{}, err
	}

	website.UpdatedAt = time.Now()
	if input.Address != "" {
		website.Address = input.Address
	}
	if input.PhoneNumber != "" {
		website.PhoneNumber = input.PhoneNumber
	}
	if input.Latitude != "" {
		website.Latitude = input.Latitude
	}
	if input.Longitude != "" {
		website.Longitude = input.Longitude
	}
	if input.Email != "" {
		website.Email = input.Email
	}
	if input.UrlInstagram != "" {
		website.URLInstagram = input.UrlInstagram
	}
	if input.UrlTikTok != "" {
		website.URLFacebook = input.UrlTikTok
	}
	if input.AboutUsBriefHomeEn != "" {
		website.AboutUsBriefHomeEn = input.AboutUsBriefHomeEn
	}
	if input.AboutUsEn != "" {
		website.AboutUsEn = input.AboutUsEn
	}
	if input.AboutUsID != "" {
		website.AboutUsID = input.AboutUsID
	}
	if input.AboutUsBriefHomeID != "" {
		website.AboutUsBriefHomeID = input.AboutUsBriefHomeID
	}
	if input.VideoWeb != "" {
		website.VideoWeb = input.VideoWeb
	}
	if input.VideoMobile != "" {
		website.VideoMobile = input.VideoMobile
	}
	if input.MetaTitle != "" {
		website.MetaTitle = input.MetaTitle
	}
	if input.MetaDesc != "" {
		website.MetaDesc = input.MetaDesc
	}
	if input.MetaKeyword != "" {
		website.MetaKeyword = input.MetaKeyword
	}
	if input.OgImage != "" {
		website.OgImage = input.OgImage
	}

	if err := tx.Save(&website).Error; err != nil {
		tx.Rollback()
		return models.Website{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.Website{}, err
	}

	return website, nil
}
