package controllers

import (
	"net/http"

	"github.com/charis16/luminor-golang-be/models"
	"github.com/charis16/luminor-golang-be/services"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users := services.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := services.CreateUser(input)
	c.JSON(http.StatusCreated, user)
}
