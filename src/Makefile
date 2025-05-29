# 📦 Makefile untuk development backend Go
GO := go
MIGRATION_DIR := ./migrations

# ⬇️ Load variable dari file .env
include .env
export

# 🆘 Help menu
help: ## Menampilkan daftar perintah yang tersedia
	@echo "📦 Daftar perintah:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 🔁 Jalankan Air untuk auto-reload saat dev
dev: ## Jalankan aplikasi dengan live reload (Air)
	air

# ▶️ Jalankan aplikasi langsung
run: ## Jalankan aplikasi secara langsung
	$(GO) run main.go

# 🧱 Buat file migration baru (interaktif)
new-migration: ## Buat file migration baru
	@read -p "Migration name (snake_case): " name; \
	mkdir -p $(MIGRATION_DIR); \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	touch "$(MIGRATION_DIR)/$${timestamp}_$${name}.up.sql"; \
	touch "$(MIGRATION_DIR)/$${timestamp}_$${name}.down.sql"; \
	echo "✅ Created $(MIGRATION_DIR)/$${timestamp}_$${name}.up.sql & down.sql"

# 🚀 Jalankan semua migration (UP)
migrate-up: ## Jalankan migration up (local)
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-up-docker: ## Jalankan migration up via Docker
	docker run --rm \
		--network shared-net \
		-v $(PWD)/$(MIGRATION_DIR):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" \
		up

# ⏪ Rollback 1 langkah (DOWN)
migrate-down: ## Rollback 1 step migration (local)
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

migrate-down-docker: ## Rollback migration 1 step via Docker
	docker run --rm \
		--network shared-net \
		-v $(PWD)/$(MIGRATION_DIR):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" \
		down 1

# 🧱 Generate model dari DB (via ./cmd/gen)
gen-model: ## Generate model GORM dari database
	$(GO) run ./cmd/gen

# 🔄 Bersihkan file hasil generate model
clean-models: ## Hapus file *_gen.go
	rm -f models/*_gen.go

regen-models: clean-models gen-model ## Bersihkan dan generate ulang model

# 📥 Install & bersihkan dependency
deps: ## Jalankan go mod tidy
	$(GO) mod tidy

# 🎯 Format semua file Go
fmt: ## Jalankan go fmt
	$(GO) fmt ./...

# 🌱 Seed data awal
seed: ## Jalankan seeder
	$(GO) run ./cmd/seeder

seed-docker: ## Jalankan seeder di dalam container
	docker exec -it luminor-api ./seeder