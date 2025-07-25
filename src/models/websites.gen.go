package models

import (
	"time"
)

const TableNameWebsite = "websites"

// Website mapped from table <websites>
type Website struct {
	ID                 int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AboutUsBriefHomeEn string    `gorm:"column:about_us_brief_home_en" json:"about_us_brief_home_en"`
	AboutUsEn          string    `gorm:"column:about_us_en" json:"about_us_en"`
	AboutUsID          string    `gorm:"column:about_us_id" json:"about_us_id"`
	AboutUsBriefHomeID string    `gorm:"column:about_us_brief_home_id" json:"about_us_brief_home_id"`
	Address            string    `gorm:"column:address" json:"address"`
	PhoneNumber        string    `gorm:"column:phone_number" json:"phone_number"`
	Email              string    `gorm:"column:email" json:"email"`
	URLInstagram       string    `gorm:"column:url_instagram" json:"url_instagram"`
	URLFacebook        string    `gorm:"column:url_facebook" json:"url_facebook"`
	VideoWeb           string    `gorm:"column:video_web" json:"video_web"`
	VideoMobile        string    `gorm:"column:video_mobile" json:"video_mobile"`
	MetaTitle          string    `gorm:"column:meta_title" json:"meta_title"`
	MetaDesc           string    `gorm:"column:meta_desc" json:"meta_desc"`
	MetaKeyword        string    `gorm:"column:meta_keyword" json:"meta_keyword"`
	OgImage            string    `gorm:"column:og_image" json:"og_image"`
	CreatedAt          time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UUID               string    `gorm:"column:uuid;default:gen_random_uuid()" json:"uuid"`
	URLTiktok          string    `gorm:"column:url_tiktok" json:"url_tiktok"`
}

// TableName Website's table name
func (*Website) TableName() string {
	return TableNameWebsite
}
