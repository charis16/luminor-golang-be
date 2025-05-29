#!/bin/sh

set -e

echo "🔧 Starting build and deployment process..."

echo "📦 Stopping existing containers..."
docker-compose down

echo "📥 Pulling latest images..."
docker-compose pull

echo "🚀 Building and starting containers..."
docker-compose up --build -d

echo "⏳ Waiting for Postgres to be ready..."
# Tunggu sampai postgres bisa diakses (timeout 30s)
RETRIES=30
until docker exec shared-postgres pg_isready -U "$POSTGRES_USER" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "⏳ Waiting... ($RETRIES)"
  sleep 1
  RETRIES=$((RETRIES - 1))
done

if [ $RETRIES -eq 0 ]; then
  echo "❌ Postgres did not become ready in time."
  exit 1
fi

echo "🧪 Checking if 'luminor' database exists..."
if ! docker exec shared-postgres psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='luminor'" | grep -q 1; then
  echo "🆕 Creating database 'luminor'..."
  docker exec shared-postgres psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE luminor"
else
  echo "✅ Database 'luminor' already exists."
fi

echo "✅ Build and run process completed"