package config

import (
	"log"

	"github.com/charis16/luminor-golang-be/src/utils"
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
	dsn := "host=" + utils.GetEnvOrPanic("DB_HOST") +
		" user=" + utils.GetEnvOrPanic("DB_USER") +
		" password=" + utils.GetEnvOrPanic("DB_PASSWORD") +
		" dbname=" + utils.GetEnvOrPanic("DB_NAME") +
		" port=" + utils.GetEnvOrPanic("DB_PORT") +
		" sslmode=" + utils.GetEnvOrPanic("DB_SSLMODE")

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
