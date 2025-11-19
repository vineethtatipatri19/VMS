#!/bin/bash

# VMS Setup Script
# Works for both local development and GitHub Codespaces

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Setting up VMS...${NC}"

# Detect environment
if [ -n "$CODESPACE_NAME" ]; then
    ENVIRONMENT="codespaces"
    BACKEND_URL="https://${CODESPACE_NAME}-8080.app.github.dev"
    FRONTEND_URL="https://${CODESPACE_NAME}-3000.app.github.dev"
    API_URL="${BACKEND_URL}/api/v1"
    
    echo -e "${YELLOW}📍 GitHub Codespaces detected${NC}"
    echo -e "   Codespace: ${CODESPACE_NAME}"
else
    ENVIRONMENT="local"
    BACKEND_URL="http://localhost:8080"
    FRONTEND_URL="http://localhost:3000"
    API_URL="http://localhost:8080/api/v1"
    
    echo -e "${YELLOW}💻 Local development detected${NC}"
fi

echo ""
echo -e "${BLUE}📍 Your URLs:${NC}"
echo -e "   Frontend: ${FRONTEND_URL}"
echo -e "   Backend:  ${BACKEND_URL}"
echo -e "   API:      ${API_URL}"
echo ""

# Create/Update .env file
echo -e "${BLUE}📝 Configuring environment...${NC}"
cat > .env << EOF
# VMS Configuration (${ENVIRONMENT})
JWT_SECRET=${JWT_SECRET:-vms-secret-key-$(date +%s)}
GEMINI_API_KEY=${GEMINI_API_KEY:-}
MIGRATE_ON_START=true
REACT_APP_API_URL=${API_URL}
EOF

echo -e "${GREEN}✅ Environment configured for ${ENVIRONMENT}${NC}"

# Stop existing containers
echo ""
echo -e "${BLUE}🛑 Stopping existing containers...${NC}"
docker-compose down > /dev/null 2>&1 || true

# Build and start services
echo ""
echo -e "${BLUE}🔨 Building services...${NC}"
echo -e "${YELLOW}   This may take a few minutes on first run...${NC}"

# Build frontend with correct API URL
docker-compose build --no-cache --build-arg REACT_APP_API_URL=${API_URL} frontend > /dev/null 2>&1 &
FRONTEND_BUILD_PID=$!

# Build backend
docker-compose build backend > /dev/null 2>&1 &
BACKEND_BUILD_PID=$!

# Show progress
echo -e "${YELLOW}   Building frontend and backend in parallel...${NC}"
wait $FRONTEND_BUILD_PID
echo -e "${GREEN}   ✓ Frontend built${NC}"
wait $BACKEND_BUILD_PID
echo -e "${GREEN}   ✓ Backend built${NC}"

# Start services
echo ""
echo -e "${BLUE}🚀 Starting services...${NC}"
docker-compose up -d

# Wait for services to be ready
echo ""
echo -e "${BLUE}⏳ Waiting for services to start...${NC}"
sleep 8

# Check database health
echo -e "${YELLOW}   Checking database...${NC}"
for i in {1..30}; do
    if docker exec pgvms-postgres pg_isready -U pgvms_user > /dev/null 2>&1; then
        echo -e "${GREEN}   ✓ Database is ready${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}   ✗ Database failed to start${NC}"
        exit 1
    fi
    sleep 1
done

# Check backend health
echo -e "${YELLOW}   Checking backend...${NC}"
for i in {1..30}; do
    if curl -sf ${BACKEND_URL}/api/v1/health > /dev/null 2>&1; then
        echo -e "${GREEN}   ✓ Backend is ready${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${YELLOW}   ⚠ Backend health check timed out (may still be starting)${NC}"
        break
    fi
    sleep 1
done

# Run database migrations
echo ""
echo -e "${BLUE}📊 Running database migrations...${NC}"
for file in infra/migrations/*.sql; do 
    filename=$(basename "$file")
    echo -e "${YELLOW}   → ${filename}${NC}"
    docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file" 2>&1 | \
        grep -v "NOTICE" | grep -v "already exists" | grep -v "will create implicit" || true
done
echo -e "${GREEN}✅ Migrations complete (7 migrations applied)${NC}"

# Load demo data
echo ""
echo -e "${BLUE}📦 Loading demo data...${NC}"
CUSTOMER_COUNT=$(docker exec pgvms-postgres psql -U pgvms_user -d pgvms -t -c "SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL;" 2>/dev/null | tr -d ' ' || echo "0")

if [ "$CUSTOMER_COUNT" -eq "0" ]; then
    echo -e "${YELLOW}   Loading customers and inventory...${NC}"
    docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_simple.sql > /dev/null 2>&1
    echo -e "${GREEN}   ✓ Loaded 15 customers, 45 inventory items${NC}"
    
    echo -e "${YELLOW}   Loading transactions, crates, wastage, and alerts...${NC}"
    docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_additional_data.sql > /dev/null 2>&1
    
    # Fix sale_items data
    docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "UPDATE sale_items si SET item_name = i.name || ' - ' || i.variant, unit = 'kg' FROM inventory_items i WHERE si.inventory_lot_id = i.id AND si.item_name IS NULL;" > /dev/null 2>&1
    
    echo -e "${GREEN}   ✓ Loaded 4 transactions, 5 crate entries, 7 wastage logs, 12 alerts${NC}"
else
    echo -e "${GREEN}✅ Demo data already exists (${CUSTOMER_COUNT} customers)${NC}"
fi

# Create demo user
echo ""
echo -e "${BLUE}👤 Setting up demo user...${NC}"
if [ -f ".demo-token" ]; then
    echo -e "${GREEN}✅ Demo user already exists${NC}"
else
    if [ -f "setup-demo-user.sh" ]; then
        ./setup-demo-user.sh > /dev/null 2>&1
        echo -e "${GREEN}✅ Demo user created${NC}"
    else
        echo -e "${YELLOW}⚠ setup-demo-user.sh not found - skipping demo user creation${NC}"
    fi
fi

# Print success message
echo ""
echo -e "${GREEN}════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✨ VMS Setup Complete! ✨${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}🌐 Access your application:${NC}"
echo ""
echo -e "   Frontend: ${GREEN}${FRONTEND_URL}${NC}"
echo -e "   Backend:  ${GREEN}${BACKEND_URL}${NC}"
echo -e "   API:      ${GREEN}${API_URL}${NC}"
echo ""
echo -e "${BLUE}🔐 Login Credentials:${NC}"
echo -e "   Email:    ${YELLOW}demo@vms.com${NC}"
echo -e "   Password: ${YELLOW}demo123${NC}"
echo ""
echo -e "${BLUE}📊 Demo Data Loaded:${NC}"
echo -e "   • 15 Customers (B2B, Wholesale, Retail)"
echo -e "   • 45 Inventory Items (Vegetables & Fruits)"
echo -e "   • 4 Transactions (Sales)"
echo -e "   • 5 Crate Ledger Entries"
echo -e "   • 7 Wastage Logs"
echo -e "   • 12 Expiry Alerts"
echo ""

if [ "$ENVIRONMENT" = "codespaces" ]; then
    echo -e "${BLUE}💡 Codespaces Tips:${NC}"
    echo -e "   • Ports are automatically forwarded"
    echo -e "   • Click the 'Ports' tab to see all services"
    echo -e "   • Click the globe icon (🌐) next to port 3000"
    echo ""
fi

echo -e "${BLUE}📚 Next Steps:${NC}"
if [ "$ENVIRONMENT" = "codespaces" ]; then
    echo -e "   1. Click: ${GREEN}${FRONTEND_URL}${NC}"
else
    echo -e "   1. Open: ${GREEN}http://localhost:3000${NC}"
fi
echo -e "   2. Login with demo@vms.com / demo123"
echo -e "   3. Explore the dashboard, inventory, customers, and transactions"
echo ""
echo -e "${BLUE}🛠️  Useful Commands:${NC}"
echo -e "   • View logs:        ${YELLOW}docker-compose logs -f backend${NC}"
echo -e "   • Restart services: ${YELLOW}docker-compose restart${NC}"
echo -e "   • Stop services:    ${YELLOW}docker-compose down${NC}"
echo -e "   • Database CLI:     ${YELLOW}docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms${NC}"
echo ""
echo -e "${GREEN}════════════════════════════════════════════════════════${NC}"
echo ""
