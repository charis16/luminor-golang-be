package main

import (
	"log"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/routes"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env not found, using default PORT 8080")
	}

	utils.InitMinio()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.1.16:3000"}, // frontend kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
	}))

	config.ConnectDB()
	v1 := r.Group("/v1/api")
	routes.RegisterUserRoutes(v1)
	routes.RegisterAuthRoutes(v1)

	port := utils.GetEnvOrDefault("PORT", "8080")

	r.Run(":" + port)
}
