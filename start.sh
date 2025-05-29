#!/bin/sh

set -e

echo "🔧 Starting build and deployment process..."

echo "📦 Stopping existing containers..."
docker-compose down

echo "📥 Pulling latest images..."
docker-compose pull

echo "🚀 Building and starting containers..."
docker-compose up --build -d

echo "⏳ Waiting for Postgres container to be ready..."
POSTGRES_USER_IN_CONTAINER=$(docker exec shared-postgres printenv POSTGRES_USER)

if [ -z "$POSTGRES_USER_IN_CONTAINER" ]; then
  echo "❌ POSTGRES_USER not found in container."
  exit 1
fi

until docker exec shared-postgres pg_isready -U "$POSTGRES_USER_IN_CONTAINER" > /dev/null 2>&1; do
  printf "."
  sleep 1
done
echo ""

echo "🧪 Checking if 'luminor' database exists..."
if ! docker exec shared-postgres psql -U "$POSTGRES_USER_IN_CONTAINER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='luminor'" | grep -q 1; then
  echo "🆕 Creating database 'luminor'..."
  docker exec shared-postgres psql -U "$POSTGRES_USER_IN_CONTAINER" -d postgres -c "CREATE DATABASE luminor"
else
  echo "✅ Database 'luminor' already exists."
fi

echo "✅ Build and run process completed."