package dto

import "time"

type UserResponse struct {
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Photo        string    `json:"photo"`
	Description  string    `json:"description"`
	Role         string    `json:"role"`
	Slug         string    `json:"slug"`
	PhoneNumber  string    `json:"phone_number"`
	URLInstagram string    `json:"url_instagram"`
	URLTikTok    string    `json:"url_tiktok"`
	URLFacebook  string    `json:"url_facebook"`
	URLYoutube   string    `json:"url_youtube"`
	IsPublished  bool      `json:"is_published"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserPortfolioResponse struct {
	User       UserResponse       `json:"users"`
	Categories []CategoryResponse `json:"categories"`
}
