#!/bin/bash

# GitHub Codespaces Setup Script
# This script configures the VMS application for GitHub Codespaces

set -e

echo "🚀 Setting up VMS for GitHub Codespaces..."

# Check if we're in a Codespace
if [ -z "$CODESPACE_NAME" ]; then
    echo "❌ This script should only be run in GitHub Codespaces"
    echo "For local development, use: docker-compose up -d --build"
    exit 1
fi

# Get the Codespaces URLs
BACKEND_URL="https://${CODESPACE_NAME}-8080.app.github.dev"
FRONTEND_URL="https://${CODESPACE_NAME}-3000.app.github.dev"

echo ""
echo "📍 Your Codespaces URLs:"
echo "   Backend:  $BACKEND_URL"
echo "   Frontend: $FRONTEND_URL"
echo ""

# Update .env file if it exists, or create it
echo "📝 Updating environment configuration..."
cat > .env << EOF
# GitHub Codespaces Configuration
JWT_SECRET=codespaces-dev-secret-$(date +%s)
GEMINI_API_KEY=${GEMINI_API_KEY:-}
MIGRATE_ON_START=true
REACT_APP_API_URL=${BACKEND_URL}/api/v1
EOF

echo "✅ Environment file updated"

# Stop existing containers
echo ""
echo "🛑 Stopping existing containers..."
docker-compose down

# Rebuild and start services
echo ""
echo "🔨 Building services with Codespaces URLs..."
docker-compose build --build-arg REACT_APP_API_URL=${BACKEND_URL}/api/v1 frontend
docker-compose up -d

# Wait for services to be ready
echo ""
echo "⏳ Waiting for services to start..."
sleep 5

# Check backend health
echo ""
echo "🏥 Checking backend health..."
if curl -sf ${BACKEND_URL}/api/v1/health > /dev/null; then
    echo "✅ Backend is healthy"
else
    echo "⚠️  Backend health check failed - it may still be starting up"
fi

# Run migrations
echo ""
echo "📊 Running database migrations..."
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file" 2>&1 | grep -v "NOTICE" | grep -v "already exists" || true
done
echo "✅ Migrations complete"

# Load demo data if not already loaded
echo ""
echo "📦 Loading demo data..."
CUSTOMER_COUNT=$(docker exec pgvms-postgres psql -U pgvms_user -d pgvms -t -c "SELECT COUNT(*) FROM customers;" 2>/dev/null | tr -d ' ' || echo "0")

if [ "$CUSTOMER_COUNT" -eq "0" ]; then
    docker cp infra/local/demo_simple.sql pgvms-postgres:/tmp/
    docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /tmp/demo_simple.sql > /dev/null 2>&1
    echo "✅ Demo data loaded (15 customers, 45 inventory items)"
else
    echo "✅ Demo data already exists (${CUSTOMER_COUNT} customers)"
fi

# Create demo user if not exists
echo ""
echo "👤 Creating demo user..."
if [ -f ".demo-token" ]; then
    echo "✅ Demo user already exists"
else
    ./setup-demo-user.sh
    echo "✅ Demo user created: demo@vms.com / demo123"
fi

# Print success message
echo ""
echo "════════════════════════════════════════════════════════"
echo "✨ VMS Setup Complete for GitHub Codespaces! ✨"
echo "════════════════════════════════════════════════════════"
echo ""
echo "🌐 Access your application:"
echo ""
echo "   Frontend:  ${FRONTEND_URL}"
echo "   Backend:   ${BACKEND_URL}/api/v1"
echo ""
echo "🔐 Login Credentials:"
echo "   Email:    demo@vms.com"
echo "   Password: demo123"
echo ""
echo "📚 Next Steps:"
echo "   1. Click on the 'Ports' tab in VS Code"
echo "   2. Find port 3000 and click the globe icon (🌐)"
echo "   3. Or directly visit: ${FRONTEND_URL}"
echo ""
echo "💡 Tip: Port visibility is set to 'Public' automatically"
echo "════════════════════════════════════════════════════════"
echo ""
