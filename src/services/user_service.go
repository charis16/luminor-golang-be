package services

import (
	"github.com/charis16/luminor-golang-be/config"
	"github.com/charis16/luminor-golang-be/models"
)

func GetAllUsers() []models.User {
	var users []models.User
	config.DB.Find(&users)
	return users
}

func CreateUser(user models.User) models.User {
	config.DB.Create(&user)
	return user
}
