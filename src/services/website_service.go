package services

import (
	"fmt"
	"time"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/dto"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
)

type WebsiteInput struct {
	Address            string `json:"address,omitempty"`
	PhoneNumber        string `json:"phone_number,omitempty"`
	Email              string `json:"email,omitempty"`
	UrlInstagram       string `json:"url_instagram,omitempty"`
	UrlTikTok          string `json:"url_tiktok,omitempty"`
	AboutUsBriefHomeEn string `json:"about_us_brief_home_en,omitempty"`
	AboutUsEn          string `json:"about_us_en,omitempty"`
	AboutUsID          string `json:"about_us_id,omitempty"`
	AboutUsBriefHomeID string `json:"about_us_brief_home_id,omitempty"`
	VideoWeb           string `json:"video_web,omitempty"`
	VideoMobile        string `json:"video_mobile,omitempty"`
	MetaTitle          string `json:"meta_title,omitempty"`
	MetaDesc           string `json:"meta_desc,omitempty"`
	MetaKeyword        string `json:"meta_keyword,omitempty"`
	OgImage            string `json:"og_image,omitempty"`
}

func GetWebsite() (dto.WebsiteResponse, int64, error) {
	var website models.Website
	if err := config.DB.First(&website).Error; err != nil {
		return dto.WebsiteResponse{}, 0, err
	}

	response := dto.WebsiteResponse{
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
		Email:              website.Email,
		UrlInstagram:       website.URLInstagram,
		UrlTikTok:          website.URLFacebook,
		CreatedAt:          website.CreatedAt,
		UpdatedAt:          website.UpdatedAt,
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

func GetWebsiteByUUID(uuid string) (models.Website, error) {
	var website models.Website
	if err := config.DB.Where("uuid = ?", uuid).First(&website).Error; err != nil {
		return models.Website{}, err
	}
	return website, nil
}

func DeleteWebsiteInformation(data models.Website, status string) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	if status == "video_web" {
		if err := utils.DeleteFromMinio("websites", data.VideoWeb); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete websites video_web photo: %v", err)
		}

		data.VideoWeb = ""
	} else if status == "video_mobile" {
		if err := utils.DeleteFromMinio("websites", data.VideoMobile); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete websites video_mobile photo: %v", err)
		}
		data.VideoMobile = ""
	} else if status == "og_image" {
		if err := utils.DeleteFromMinio("websites", data.OgImage); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete websites og_image photo: %v", err)
		}
		data.OgImage = ""
	}

	if err := tx.Save(&data).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update website information: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
