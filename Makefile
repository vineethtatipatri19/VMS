
### PGVMS Makefile (updated)

BACKEND_DIR=backend
FRONTEND_DIR=frontend
LOCAL_INFRA_DIR=infra/local

### Load env from backend/.env if present (not exported here)
export DATABASE_URL=$(shell [ -f backend/.env ] && sed -n 's/^DATABASE_URL=\(.*\)/\1/p' backend/.env)

dev:
	@echo "🚀 Starting PGVMS (local backend + local DB + frontend)…"
	cd $(LOCAL_INFRA_DIR) && docker compose up --build -d
	cd $(BACKEND_DIR) && air &
	cd $(FRONTEND_DIR) && npm start

run:
	@echo "▶ Running backend locally with Go"
	cd $(BACKEND_DIR) && go run main.go

test:
	cd $(BACKEND_DIR) && go test ./...

docker-build:
	cd $(BACKEND_DIR) && docker build -t pgvms-backend .

docker-run:
	docker run -p 8080:8080 --env-file backend/.env pgvms-backend

db-up:
	cd $(LOCAL_INFRA_DIR) && docker compose up -d db

db-down:
	cd $(LOCAL_INFRA_DIR) && docker compose down

migrate-up:
	migrate -path infra/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path infra/migrations -database "$(DATABASE_URL)" down 1

seed:
	psql "$(DATABASE_URL)" -f infra/local/seed.sql

clean:
	go clean
	rm -f backend/pgvms-api

deploy:
	gcloud builds submit --config infra/cloudbuild.yaml .

help:
	@echo ""
	@echo "PGVMS Automation Commands"
	@echo "--------------------------"
	@echo "make dev           – Start backend + DB + frontend (BEST OPTION)"
	@echo "make run           – Run backend only"
	@echo "make test          – Run backend tests"
	@echo "make docker-build  – Build backend Docker image"
	@echo "make migrate-up    – Run DB migrations"
	@echo "make seed          – Insert sample data"
	@echo "make deploy        – Deploy backend to Google Cloud Run"
	@echo ""
