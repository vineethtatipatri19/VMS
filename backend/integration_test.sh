
#!/usr/bin/env bash
set -e
ROOT=$(cd $(dirname $0)/.. && pwd)
COMPOSE_DIR=$ROOT/infra/local
echo "Starting local infra..."
cd $COMPOSE_DIR
docker compose up -d db
# wait for postgres
for i in {1..30}; do
  docker run --rm --network host postgres:15 pg_isready -h localhost -p 5432 && break || sleep 1
done
export DATABASE_URL=postgres://pgvms:pgvms123@localhost:5432/pgvms?sslmode=disable
cd $ROOT/backend
echo "Running migrations..."
migrate -path ../infra/migrations -database "$DATABASE_URL" up || true
echo "Running go tests..."
go test ./... -v
echo "Tests done"
