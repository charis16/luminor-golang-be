package dto

import "time"

type CategoryResponse struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	PhotoUrl    string    `json:"photo_url"`
	YoutubeURL  string    `json:"youtube_url"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryBySlugResponse struct {
	UUID        string         `json:"uuid"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Slug        string         `json:"slug"`
	PhotoUrl    string         `json:"photo_url"`
	Users       []UserResponse `json:"users"`
}
