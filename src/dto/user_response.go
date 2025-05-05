package dto

import "github.com/charis16/luminor-golang-be/src/models"

type UserResponse struct {
	models.User
	ID any `json:"-"`
}
