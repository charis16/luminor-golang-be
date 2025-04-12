package main

import (
	"github.com/charis16/luminor-golang-be/config"
	"github.com/charis16/luminor-golang-be/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDB()
	routes.RegisterUserRoutes(r)
	r.Run(":8080")
}
