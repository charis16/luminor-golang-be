package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GenerateModels() {
	// ✅ Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file tidak ditemukan, fallback ke os env")
	}

	// ✅ Ambil dari env
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	// ✅ Konfig generator
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./models", // relatif dari CWD (misal: `src/models`)
		ModelPkgPath: "models",   // RELATIF terhadap nama module `go.mod`
	})

	// ✅ Buka koneksi DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ gagal konek ke DB: %v", err)
	}

	// ✅ Generate semua table
	g.UseDB(db)
	g.GenerateAllTable()
	g.Execute()
}
