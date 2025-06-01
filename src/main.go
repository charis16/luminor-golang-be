package main

import (
	"log"
	"strings"

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

	utils.InitR2()
	ginMode := utils.GetEnvOrDefault("GIN_MODE", "development")
	if ginMode != "" && ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	feUrls := utils.GetEnvOrDefault("FE_URL", "http://localhost:3000")
	originList := strings.Split(feUrls, ",")
	allowOrigins := make([]string, 0, len(originList))

	for _, origin := range originList {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			allowOrigins = append(allowOrigins, trimmed)
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins, // frontend kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
	}))

	config.ConnectDB()
	v1 := r.Group("/v1/api")
	routes.UserRoutes(v1)
	routes.AuthRoutes(v1)
	routes.FaqRoutes(v1)
	routes.CategoryRoutes(v1)
	routes.WebsiteRoutes(v1)
	routes.AlbumRoutes(v1)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := utils.GetEnvOrDefault("PORT", "8080")
	env := utils.GetEnvOrDefault("APP_ENV", "development")
	if env == "production" {
		log.Println("🚀 Running in PRODUCTION mode (no TLS)...")
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Run failed:", err)
		}
	} else {
		log.Println("🔐 Running in DEV mode with TLS...")
		if err := r.RunTLS(":"+port, "../../certs/localhost.pem", "../../certs/localhost-key.pem"); err != nil {
			log.Fatal("RunTLS failed:", err)
		}
	}
}
