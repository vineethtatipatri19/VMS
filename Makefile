
### PGVMS Makefile

BACKEND_DIR=backend
FRONTEND_DIR=frontend

### Load env from .env if present
-include .env
export

.PHONY: help dev up down build test clean

help:
	@echo ""
	@echo "PGVMS - Perishable Goods Vendor Management System"
	@echo "================================================"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make up          - Start all services with Docker Compose"
	@echo "  make down        - Stop all services"
	@echo "  make logs        - View logs from all services"
	@echo ""
	@echo "🔧 Development:"
	@echo "  make dev         - Run backend locally (requires PostgreSQL)"
	@echo "  make test        - Run backend tests"
	@echo "  make build       - Build backend binary"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Start with Docker Compose"
	@echo "  make docker-down  - Stop Docker Compose"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  make db-migrate   - Run database migrations"
	@echo "  make db-seed      - Seed sample data"
	@echo ""
	@echo "☁️  Deployment:"
	@echo "  make deploy       - Deploy to Google Cloud Run"
	@echo ""

# Quick start with Docker Compose
up: docker-up

down: docker-down

# Docker Compose commands
docker-up:
	@echo "🚀 Starting PGVMS with Docker Compose..."
	docker-compose up --build -d
	@echo "✅ Services started!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080"
	@echo ""
	@echo "Run 'make logs' to view logs"

docker-down:
	@echo "🛑 Stopping services..."
	docker-compose down
	@echo "✅ Services stopped"

docker-build:
	@echo "🔨 Building Docker images..."
	docker-compose build

logs:
	docker-compose logs -f

# Local development
dev:
	@echo "▶ Running backend locally"
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "❌ DATABASE_URL not set. Please create .env file from .env.example"; \
		exit 1; \
	fi
	cd $(BACKEND_DIR) && go run .

build:
	@echo "🔨 Building backend..."
	cd $(BACKEND_DIR) && go build -o pgvms .
	@echo "✅ Built: backend/pgvms"

test:
	@echo "🧪 Running tests..."
	cd $(BACKEND_DIR) && go test ./...

# Database
db-migrate:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "❌ DATABASE_URL not set"; \
		exit 1; \
	fi
	migrate -path infra/migrations -database "$(DATABASE_URL)" up

db-seed:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "❌ DATABASE_URL not set"; \
		exit 1; \
	fi
	psql "$(DATABASE_URL)" -f infra/local/seed.sql

# Cleanup
clean:
	@echo "🧹 Cleaning up..."
	rm -f backend/pgvms backend/pgvms-api
	rm -rf frontend/node_modules frontend/build
	docker-compose down -v
	@echo "✅ Cleanup complete"

# Deployment
deploy:
	@echo "☁️  Deploying to Google Cloud Run..."
	gcloud builds submit --config infra/cloudbuild.yaml .

