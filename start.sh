#!/bin/sh

echo "🔧 Building and running production containers..."

# Matikan dan tarik ulang container
docker-compose down
docker-compose pull
docker-compose up --build -d

echo "✅ Build and run completed."

echo "⏳ Waiting for Postgres to be ready..."
sleep 5

# Ambil nama user dari dalam container
POSTGRES_USER_IN_CONTAINER=$(docker exec shared-postgres printenv POSTGRES_USER)

echo "🧪 Checking if 'luminor' database exists..."
if ! docker exec shared-postgres psql -U "$POSTGRES_USER_IN_CONTAINER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='luminor'" | grep -q 1; then
  echo "🆕 Creating database 'luminor'..."
  docker exec shared-postgres psql -U "$POSTGRES_USER_IN_CONTAINER" -d postgres -c "CREATE DATABASE luminor"
else
  echo "✅ Database 'luminor' already exists."
fi

echo "🚀 Application is ready at http://localhost"