#!/bin/bash

# RAG Financial News Worker - Startup Script
# This script starts the complete RAG infrastructure using Docker Compose

set -e

echo "🚀 Starting RAG Financial News Worker"
echo "======================================"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p logs
mkdir -p data

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check service status
echo "📊 Checking service status..."
docker-compose ps

# Show service URLs
echo ""
echo "🌐 Service URLs:"
echo "  - Kafka UI: http://localhost:8081"
echo "  - Spark Master: http://localhost:8080"
echo "  - ChromaDB: http://localhost:8000"
echo ""

# Show logs
echo "📋 Recent logs:"
docker-compose logs --tail=20

echo ""
echo "✅ RAG Financial News Worker is starting up!"
echo "📝 Check logs with: docker-compose logs -f"
echo "🛑 Stop with: docker-compose down"
