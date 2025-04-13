# Makefile untuk development backend Go

GO := go

# ⬇️ Load variable dari file .env
include .env
export

# 🔁 Jalankan Air untuk auto-reload saat dev
dev:
	air

# ▶️ Jalankan aplikasi langsung
run:
	$(GO) run main.go

# 🛠 Generate model GORM dari database (via ./cmd/gen)
gen-model:
	$(GO) run ./cmd/gen

# 🔄 Bersihkan file hasil generate GORM
clean-models:
	rm -f models/*_gen.go

# 🧱 Buat file migration baru (interaktif)
new-migration:
	@read -p "Migration name (snake_case): " name; \
	mkdir -p migrations; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	touch "migrations/$${timestamp}_$${name}.up.sql"; \
	touch "migrations/$${timestamp}_$${name}.down.sql"; \
	echo "✅ Created migrations/$${timestamp}_$${name}.up.sql & down.sql"

# 🚀 Jalankan semua migration (UP)
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

# ⏪ Rollback 1 step (DOWN)
migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

# 🔧 Format semua file Go
fmt:
	$(GO) fmt ./...

# 📦 Install atau bersihkan dependency
deps:
	$(GO) mod tidy

# 🔄 Regenerate model full (clean + generate)
regen-models: clean-models gen-model