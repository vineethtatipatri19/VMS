#!/bin/bash

# PGVMS Quick Start Script
echo "🚀 PGVMS - Perishable Goods Vendor Management System"
echo "===================================================="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please edit .env and set your JWT_SECRET before running again"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install it first."
    exit 1
fi

echo "✅ All prerequisites met"
echo ""
echo "🐳 Starting services with Docker Compose..."
docker-compose up --build -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

echo ""
echo "✅ PGVMS is running!"
echo ""
echo "📍 Access the application:"
echo "   Frontend:  http://localhost:3000"
echo "   Backend:   http://localhost:8080"
echo "   Health:    http://localhost:8080/api/v1/health"
echo ""
echo "📚 Next steps:"
echo "   1. Open http://localhost:3000 in your browser"
echo "   2. Click 'Register here' to create an account"
echo "   3. Login and start using the system"
echo ""
echo "🔍 View logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Stop services:"
echo "   docker-compose down"
echo ""
