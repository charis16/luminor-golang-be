package config

import (
	"fmt"
	"log"

	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, fallback ke system env")
	}

	// Susun DSN dari variabel terpisah
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		utils.GetEnvOrPanic("DB_HOST"),
		utils.GetEnvOrPanic("DB_USER"),
		utils.GetEnvOrPanic("DB_PASSWORD"),
		utils.GetEnvOrPanic("DB_NAME"),
		utils.GetEnvOrPanic("DB_PORT"),
		utils.GetEnvOrPanic("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek DB:", err)
	}

	DB = db
	fmt.Println("✅ Database connected")
}
