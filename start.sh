#!/bin/sh

echo "🔧 Building and running production containers..."

# Jalankan docker-compose (build & up)
docker-compose up --build -d

echo "✅ App is running at http://localhost"