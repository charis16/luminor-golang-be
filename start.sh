#!/bin/sh

set -e

echo "🔧 Starting build and deployment process..."

echo "📦 Stopping existing containers..."
docker-compose down --volumes --remove-orphans

echo "🧹 Cleaning unused containers and networks..."
docker container prune -f
docker volume prune -f
docker network prune -f

echo "📥 Pulling latest images..."
docker-compose pull

echo "🚀 Building and starting containers..."
docker-compose build
docker-compose up -d

echo "⏳ Waiting for Postgres container to start..."
until docker inspect -f '{{.State.Running}}' shared-postgres 2>/dev/null | grep true > /dev/null; do
  printf "."
  sleep 1
done
echo ""

echo "⏳ Waiting for Postgres to be ready..."
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