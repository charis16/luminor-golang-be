#!/bin/sh

set -e

echo "🔧 Starting build and deployment process..."

echo "📦 Stopping existing containers..."
docker-compose down

echo "📥 Pulling latest images..."
docker-compose pull

echo "🚀 Building and starting containers..."
docker-compose up --build -d

# Ambil nama user dari dalam container
POSTGRES_USER_IN_CONTAINER=$(docker exec shared-postgres printenv POSTGRES_USER)

echo "🧪 Checking if 'luminor' database exists..."
if ! docker exec shared-postgres psql -U "$POSTGRES_USER_IN_CONTAINER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='luminor'" | grep -q 1; then
  echo "🆕 Creating database 'luminor'..."
  docker exec shared-postgres psql -U "$POSTGRES_USER_IN_CONTAINER" -d postgres -c "CREATE DATABASE luminor"
else
  echo "✅ Database 'luminor' already exists."
fi

echo "✅ Build and run process completed."