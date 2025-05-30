#!/bin/sh

set -e

echo "🔧 Starting build and deployment process..."

# Step 1: Stop existing containers and clean up
echo "📦 Stopping existing containers and removing orphans..."
docker-compose down --volumes --remove-orphans

echo "🧹 Cleaning up unused Docker resources..."
docker container prune -f
docker volume prune -f
docker network prune -f
docker images -f "dangling=true" -q | xargs -r docker rmi -f

# 🔥 Hapus image dangling (<none>)
docker images -f "dangling=true" -q | xargs -r docker rmi -f

# 🔥 Hapus cache build yang tidak dipakai
docker builder prune -af

# Step 2: Pull and rebuild images
echo "📥 Pulling latest base images..."
docker-compose pull

echo "🔨 Building containers from Dockerfile..."
docker-compose build --no-cache --force-rm

# Step 3: Start containers
echo "🚀 Starting containers..."
docker-compose up -d

# Step 4: Wait for Postgres to become available
echo "⏳ Waiting for 'shared-postgres' container to start..."
until docker inspect -f '{{.State.Running}}' shared-postgres 2>/dev/null | grep true > /dev/null; do
  printf "."
  sleep 1
done
echo ""

echo "⏳ Waiting for Postgres service to be ready..."
POSTGRES_USER=$(docker exec shared-postgres printenv POSTGRES_USER)

if [ -z "$POSTGRES_USER" ]; then
  echo "❌ Environment variable POSTGRES_USER not found in container."
  exit 1
fi

until docker exec shared-postgres pg_isready -U "$POSTGRES_USER" > /dev/null 2>&1; do
  printf "."
  sleep 1
done
echo ""

# Step 5: Ensure database exists
echo "🧪 Verifying 'luminor' database..."
if ! docker exec shared-postgres psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='luminor'" | grep -q 1; then
  echo "🆕 Creating 'luminor' database..."
  docker exec shared-postgres psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE luminor"
else
  echo "✅ Database 'luminor' already exists."
fi

echo "✅ Build and deployment process completed successfully."