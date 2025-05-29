package dto

import "time"

type WebsiteResponse struct {
	UUID               string    `json:"uuid"`
	Address            string    `json:"address"`
	PhoneNumber        string    `json:"phone_number"`
	Email              string    `json:"email"`
	UrlInstagram       string    `json:"url_instagram"`
	UrlTikTok          string    `json:"url_tiktok"`
	IsPublished        bool      `json:"is_published"`
	AboutUsBriefHomeEn string    `json:"about_us_brief_home_en"`
	AboutUsEn          string    `json:"about_us_en"`
	AboutUsID          string    `json:"about_us_id"`
	AboutUsBriefHomeID string    `json:"about_us_brief_home_id"`
	VideoWeb           string    `json:"video_web"`
	VideoMobile        string    `json:"video_mobile"`
	MetaTitle          string    `json:"meta_title"`
	MetaDesc           string    `json:"meta_desc"`
	MetaKeyword        string    `json:"meta_keyword"`
	OgImage            string    `json:"og_image"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
