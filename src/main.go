package main

import (
	"log"
	"os"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env not found, using default PORT 8080")
	}

	r := gin.Default()
	config.ConnectDB()
	routes.RegisterUserRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback kalau .env tidak ada
	}

	r.Run(":" + port)
}
