package dto

import (
	"time"
)

type AlbumResponse struct {
	UUID         string    `json:"uuid"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	CategoryId   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	UserAvatar   string    `json:"user_avatar"`
	Description  string    `json:"description"`
	Thumbnail    string    `json:"thumbnail"`
	Images       []string  `json:"images"` // ubah jadi array string
	IsPublished  bool      `json:"is_published"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
