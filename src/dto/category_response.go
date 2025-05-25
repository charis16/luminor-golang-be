package dto

import "time"

type CategoryResponse struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryOption struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
