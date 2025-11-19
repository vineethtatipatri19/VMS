# VMS Quick Reference Guide

## 🚀 Quick Start

**One command setup (works on local & Codespaces):**
```bash
bash setup.sh
```

**Login:** demo@vms.com / demo123

**URLs:**
- Local: http://localhost:3000
- Codespaces: Displayed after setup completes

## Demo Data Summary

After running `setup.sh`, you have:
- **15 Customers**: 3 B2B, 2 Wholesale, 7 Retail, 2 B2C, 1 Blocked
- **45 Inventory Items**: 25 vegetables, 20 fruits (with variants)
- **4 Transactions**: Credit sale, wholesale sale, retail sales with complete sale items
- **5 Crate Entries**: Issued and returned with balance tracking
- **7 Wastage Logs**: Spoiled, damaged, expired reasons
- **12 Expiry Alerts**: Critical to medium priority alerts

## Essential Commands

### Start/Stop Services

```bash
# Complete automated setup
bash setup.sh

# Start everything manually
docker-compose up -d --build

# Stop everything
docker-compose down

# Complete cleanup (removes volumes)
docker-compose down -v

# Restart a service
docker-compose restart backend
docker-compose restart frontend

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Check status
docker-compose ps
```

### Database Operations

```bash
# Access database
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms

# Check tables
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"

# View recent transactions
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT id, date, type, total_amount FROM transactions ORDER BY date DESC LIMIT 10;"

# Count all records
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "
SELECT 
  'customers' as table, COUNT(*) FROM customers WHERE deleted_at IS NULL
UNION ALL
SELECT 'inventory', COUNT(*) FROM inventory_items WHERE deleted_at IS NULL
UNION ALL
SELECT 'transactions', COUNT(*) FROM transactions WHERE deleted_at IS NULL
UNION ALL
SELECT 'sale_items', COUNT(*) FROM sale_items
UNION ALL
SELECT 'crate_ledger', COUNT(*) FROM crate_ledger
UNION ALL
SELECT 'wastage_log', COUNT(*) FROM wastage_log
UNION ALL
SELECT 'expiry_alerts', COUNT(*) FROM expiry_alerts;
"

# Manual migration execution (if needed)
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Reload demo data
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_simple.sql
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_additional_data.sql

# Fix sale_items data
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "
UPDATE sale_items si 
SET item_name = i.name || ' - ' || i.variant, unit = 'kg' 
FROM inventory_items i 
WHERE si.inventory_lot_id = i.id AND si.item_name IS NULL;
"
```

### API Testing

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Login and get token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@vms.com","password":"demo123"}' | jq -r '.token'

# Save token for subsequent requests
TOKEN=$(cat .demo-token)

# Get all endpoints with data
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers | jq '.data | length'
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/inventory | jq '.data | length'
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/transactions | jq '.data | length'
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/crates | jq '.data | length'
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/wastage | jq '.data | length'
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/expiry-alerts | jq '.data | length'

# Get dashboard stats
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/dashboard | jq

# Get transaction with sale items
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/transactions" | jq '.data[0]'
```

### Development Workflow

```bash
# Backend changes
cd backend
go fmt ./...          # Format code
go test ./...         # Run tests
go build -o pgvms .   # Build binary
docker-compose up -d --build backend  # Rebuild container

# Frontend changes
cd frontend
npm run build         # Build production bundle
docker-compose up -d --build frontend  # Rebuild container

# Database changes
# 1. Create new migration in infra/migrations/
# 2. Restart backend with MIGRATE_ON_START=true
docker-compose restart backend
```

## Default Credentials

**Demo User:**
- Email: `demo@vms.com`
- Password: `demo123`
- Token saved in: `.demo-token` (after setup.sh)

**Database:**
- Host: `localhost` (or `pgvms-postgres` from containers)
- Port: `5432`
- Database: `pgvms`
- User: `pgvms_user`
- Password: `pgvms_password`

**URLs:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1
- Health: http://localhost:8080/api/v1/health

## Demo Data Contents

After loading both `demo_simple.sql` and `demo_additional_data.sql`:

**15 Customers:**
- 3 B2B: Fresh Mart Supermarket, Green Valley Restaurant, Organic Hub
- 2 Wholesale: Ravi Wholesale Trading, Metro Vegetables
- 7 Retail: Priya Sharma, Amit Patel, Sunita Reddy, etc.
- 2 B2C: Rajesh Kumar, Deepa Mehta
- 1 Blocked: Vijay Malhotra

**45 Inventory Items:**
- 25 Vegetables: Tomato (Local, Hybrid), Onion (White, Red), Potato, Carrot, Cabbage, Spinach, Cauliflower, Brinjal, Capsicum, etc.
- 20 Fruits: Apple (Shimla, Kashmiri), Banana (Robusta, Yelakki), Orange, Mango, Grapes, Watermelon, Papaya, Pomegranate, etc.

**4 Transactions:**
- INV-20251117-1011: Credit Sale to Fresh Mart (₹1,180)
- INV-20251117-1012: Wholesale Sale to Ravi Trading (₹10,500)
- INV-20251117-1013: Retail Sale to Priya Sharma (₹770)
- INV-20251117-1014: Retail Sale to Amit Patel (₹1,080)
- All with complete sale_items (12 items total)

**5 Crate Entries:**
- 2 Issues: 100 plastic crates, 50 wooden crates
- 2 Returns: 80 plastic returned, 30 wooden returned
- 1 Adjustment: Damaged crates adjustment
- Balance tracking: 20 plastic, 20 wooden unreturned

**7 Wastage Logs:**
- Spoiled: Tomato, Banana, Milk
- Damaged: Onion, Apple
- Expired: Paneer, Rice
- Total quantity: 111kg/units, Total value: ₹2,220

**12 Expiry Alerts:**
- Critical (1-3 days): Tomato Hybrid, Curd, Paneer
- High (4-7 days): Onion White, Banana Yelakki, Milk Full Cream
- Medium (8-14 days): Carrot, Spinach, Apple Shimla, etc.

## Common Tasks

### Add New Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Customer",
    "contact_number": "9876543210",
    "email": "new@example.com",
    "address": "123 Main St",
    "customer_type": "retail"
  }'
```

### Create Transaction with Items
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-uuid",
    "date": "2025-11-18T10:00:00Z",
    "type": "sale",
    "total_amount": 5000,
    "payment_amount": 5000,
    "items": [
      {
        "inventory_lot_id": "inventory-uuid",
        "item_name": "Tomato",
        "quantity": 10,
        "unit": "kg",
        "price_per_unit": 50,
        "total": 500
      }
    ]
  }'
```

### View Dashboard Data
Frontend: http://localhost:3000/dashboard

**Shows:**
- Total customers count
- Unreturned crates
- Outstanding balance
- Expiring soon items
- Sales trend chart (last 7 days)
- Top products chart (by quantity sold)
- Recent activity

## Troubleshooting Quick Fixes

### "Cannot connect to backend"
```bash
# Check backend is running
docker ps | grep pgvms-backend
docker logs pgvms-backend --tail=50

# Check health
curl http://localhost:8080/api/v1/health

# Restart backend
docker-compose restart backend

# Full rebuild if needed
docker-compose up -d --build backend
```

### "Database connection failed"
```bash
# Check postgres is healthy
docker ps | grep pgvms-postgres

# Check logs
docker logs pgvms-postgres --tail=50

# Test connection
docker exec pgvms-postgres pg_isready -U pgvms_user -d pgvms

# Verify credentials
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms
```

### "Frontend shows blank page"
```bash
# Check frontend logs
docker logs pgvms-frontend --tail=50

# Test port
curl -I http://localhost:3000

# Check API URL configuration
docker exec pgvms-frontend cat /usr/share/nginx/html/static/js/main.*.js | grep -o "http.*api/v1"

# Rebuild frontend with correct URL
docker-compose down frontend
docker-compose up -d --build frontend
```

### "No transactions/crates/wastage showing"
```bash
# Check if data exists
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "
SELECT 
  (SELECT COUNT(*) FROM transactions) as txn_count,
  (SELECT COUNT(*) FROM crate_ledger) as crate_count,
  (SELECT COUNT(*) FROM wastage_log) as wastage_count,
  (SELECT COUNT(*) FROM expiry_alerts) as alert_count;
"

# If zero, reload data
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_additional_data.sql

# Fix sale_items if needed
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "
UPDATE sale_items si 
SET item_name = i.name || ' - ' || i.variant, unit = 'kg' 
FROM inventory_items i 
WHERE si.inventory_lot_id = i.id AND si.item_name IS NULL;
"

# Verify API responses
TOKEN=$(cat .demo-token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/transactions | jq '.success, .count'
```

### "Charts not showing data"
```bash
# Verify sale_items exist and have proper data
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "
SELECT COUNT(*) as count, 
       COUNT(item_name) as items_with_name,
       COUNT(unit) as items_with_unit
FROM sale_items;
"

# Check transaction details in API
TOKEN=$(cat .demo-token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/transactions | jq '.data[0].Details'

# If Details is null or empty, rebuild backend
docker-compose up -d --build backend
```

### "Automatic migrations not working"
```bash
# Run migrations manually
for file in infra/migrations/*.sql; do 
  echo "Running migration: $(basename $file)"
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Verify all tables exist
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"
# Should show: customers, inventory_items, transactions, sale_items, 
#              crate_ledger, wastage_log, expiry_alerts, users
```

### "Port already in use"
```bash
# Check what's using the ports
sudo lsof -i :3000  # Frontend
sudo lsof -i :8080  # Backend
sudo lsof -i :5432  # PostgreSQL

# Stop conflicting services
sudo kill -9 <PID>

# Or change ports in docker-compose.yml
# Then restart
docker-compose down && docker-compose up -d
```

### "Permission denied on setup.sh"
```bash
# Make executable
chmod +x setup.sh

# Run with bash explicitly
bash setup.sh
```

## File Locations

**Configuration:**
- Environment: `.env` (auto-created by setup.sh)
- Docker: `docker-compose.yml`
- Frontend config: `frontend/nginx.conf`
- Backend config: `backend/internal/config/config.go`

**Code:**
- Backend entry: `backend/main.go`
- Backend routes: `backend/internal/router/router.go`
- Backend handlers: `backend/internal/handlers/*.go`
- Frontend entry: `frontend/src/App.js`
- API client: `frontend/src/services/api.js`
- Dashboard: `frontend/src/pages/Dashboard.js`

**Database:**
- Migrations: `infra/migrations/*.sql` (7 files)
- Demo data: `infra/local/demo_simple.sql` (customers, inventory)
- Additional data: `infra/local/demo_additional_data.sql` (transactions, crates, wastage, alerts)

**Setup:**
- Unified setup: `setup.sh` (recommended)
- Codespaces setup: `setup-codespaces.sh` (legacy)
- Demo user: `setup-demo-user.sh` (called by setup.sh)

**Documentation:**
- This guide: `QUICK_REFERENCE.md`
- Full setup: `SETUP.md`
- Architecture: `ARCHITECTURE.md` and `backend/ARCHITECTURE.md`
- API reference: `docs/API.md`
- Implementation: `IMPLEMENTATION_SUMMARY.md`
- GitHub issues: `GITHUB_ISSUES.md`

## Useful SQL Queries

```sql
-- View all customers with balances
SELECT name, customer_type, current_balance, credit_limit 
FROM customers 
WHERE deleted_at IS NULL 
ORDER BY current_balance DESC;

-- View expiring inventory
SELECT name, variant, quantity, unit, expiry_date, 
       (expiry_date - CURRENT_DATE) as days_remaining
FROM inventory_items 
WHERE deleted_at IS NULL 
  AND expiry_date <= CURRENT_DATE + INTERVAL '7 days'
ORDER BY expiry_date;

-- View recent transactions with items
SELECT t.id, t.date, t.type, t.total_amount, c.name as customer,
       COUNT(si.id) as item_count
FROM transactions t
LEFT JOIN customers c ON t.customer_id = c.id
LEFT JOIN sale_items si ON t.id = si.transaction_id
WHERE t.deleted_at IS NULL
GROUP BY t.id, t.date, t.type, t.total_amount, c.name
ORDER BY t.date DESC
LIMIT 10;

-- Calculate total sales by date
SELECT DATE(date) as sale_date, 
       COUNT(*) as transaction_count,
       SUM(total_amount) as total_sales
FROM transactions 
WHERE type = 'sale' 
  AND deleted_at IS NULL
  AND date >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY DATE(date)
ORDER BY sale_date;

-- Top selling products
SELECT item_name, 
       SUM(quantity) as total_quantity,
       SUM(total) as total_revenue
FROM sale_items
GROUP BY item_name
ORDER BY total_quantity DESC
LIMIT 10;

-- Customer credit utilization
SELECT name, 
       credit_limit,
       current_balance,
       ROUND((current_balance / NULLIF(credit_limit, 0)) * 100, 2) as utilization_pct
FROM customers
WHERE deleted_at IS NULL 
  AND credit_limit > 0
ORDER BY utilization_pct DESC;
```

## Performance Tips

1. **Database indexes** - Already created for common queries
2. **Pagination** - Implement for large lists (future enhancement)
3. **Caching** - Use Redis for frequently accessed data (future)
4. **Connection pooling** - Already configured in backend
5. **Chart data** - Calculated once per data load, not on every render

## Next Steps After Setup

1. ✅ **Login** with demo@vms.com / demo123
2. ✅ **Explore Dashboard** - View real charts with sales trends and top products
3. ✅ **Browse Customers** - See 15 customers with different types (B2B, Wholesale, Retail)
4. ✅ **Check Inventory** - 45 items with FEFO sorting, variants, and expiry tracking
5. ✅ **View Transactions** - 4 complete transactions with sale items details
6. ✅ **Crate Management** - Track issued/returned crates with balance
7. ✅ **Wastage Logs** - See 7 wastage entries with reasons
8. ✅ **Expiry Alerts** - 12 alerts by priority (critical/high/medium)
9. 🔨 **Create New Records** - Add customers, inventory, transactions
10. 🔨 **Generate Reports** - Sales, inventory, customer reports
11. 🔨 **Test AI Forecasting** - If GEMINI_API_KEY is set in .env
12. 🔨 **Try Editing** - Update records and see soft-delete in action
13. 🔨 **Mobile Testing** - React Native app if deployed

## Support Resources

- **Setup Guide:** [SETUP.md](SETUP.md) - Comprehensive setup instructions
- **Quick Reference:** This file (QUICK_REFERENCE.md)
- **Architecture:** [ARCHITECTURE.md](ARCHITECTURE.md), [backend/ARCHITECTURE.md](backend/ARCHITECTURE.md)
- **API Docs:** [docs/API.md](docs/API.md)
- **Implementation:** [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)
- **GitHub:** https://github.com/vineethtatipatri19/VMS

---

## 🎯 Quick Start Command (All-in-One)

```bash
# Clone, setup, and run (takes ~2-3 minutes)
git clone https://github.com/vineethtatipatri19/VMS.git && \
cd VMS && \
bash setup.sh

# After setup completes, open the URL shown and login with:
# Email: demo@vms.com
# Password: demo123
```

**That's it! Your VMS is ready with complete demo data.** 🎉
