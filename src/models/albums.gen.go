package models

import (
	"time"

	"github.com/lib/pq"
)

const TableNameAlbum = "albums"

// Album mapped from table <albums>
type Album struct {
	ID          int32          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UUID        string         `gorm:"column:uuid;default:gen_random_uuid()" json:"uuid"`
	Slug        string         `gorm:"column:slug" json:"slug"`
	Title       string         `gorm:"column:title" json:"title"`
	CategoryID  int32          `gorm:"column:category_id" json:"category_id"`
	Description string         `gorm:"column:description" json:"description"`
	YoutubeURL  string         `gorm:"column:youtube_url" json:"youtube_url"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
	Thumbnail   string         `gorm:"column:thumbnail" json:"thumbnail"`
	IsPublished bool           `gorm:"column:is_published" json:"is_published"`
	CreatedAt   time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UserID      int32          `gorm:"column:user_id" json:"user_id"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
}

// TableName Album's table name
func (*Album) TableName() string {
	return TableNameAlbum
}
