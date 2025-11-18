# VMS Quick Reference Guide

## Essential Commands

### Start/Stop Services

```bash
# Start everything
docker-compose up -d --build

# Stop everything
docker-compose down

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

# Load demo data
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /docker-entrypoint-initdb.d/demo_simple.sql

# Check tables
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"

# View transactions
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT id, date, type, total_amount FROM transactions ORDER BY date DESC LIMIT 10;"

# Count records
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "
SELECT 
  'customers' as table, COUNT(*) as count FROM customers WHERE deleted_at IS NULL
UNION ALL
SELECT 'inventory', COUNT(*) FROM inventory_items WHERE deleted_at IS NULL
UNION ALL
SELECT 'transactions', COUNT(*) FROM transactions WHERE deleted_at IS NULL
UNION ALL
SELECT 'sale_items', COUNT(*) FROM sale_items;
"
```

### API Testing

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Register user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@vms.com","password":"demo123"}'

# Get customers (requires token)
TOKEN="your-jwt-token"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers

# Get dashboard stats
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/dashboard

# Get transactions with sale items
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/transactions | jq '.data[0].Details'
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

**Demo User (after running setup-demo-user.sh):**
- Email: `demo@vms.com`
- Password: `demo123`

**Database:**
- Host: `localhost`
- Port: `5432`
- Database: `pgvms`
- User: `pgvms_user`
- Password: `pgvms_password`

**URLs:**
- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- API: http://localhost:8080/api/v1

## Demo Data Contents

After loading `demo_simple.sql`:

**15 Customers:**
- 3 B2B (Fresh Mart Supermarket, Green Valley Restaurant, Organic Hub)
- 2 Wholesale (Ravi Wholesale Trading, Metro Vegetables)
- 10 Retail (Individual customers)

**45 Inventory Items:**
- Vegetables (20): Tomato, Onion, Potato, Carrot, Cabbage, Spinach, Cauliflower, etc.
- Fruits (10): Apple, Banana, Orange, Mango, Grapes, Watermelon, etc.
- Dairy (8): Milk, Paneer, Curd, Butter, Cheese, Ghee, etc.
- Grains (7): Rice, Wheat, Lentils, Chickpeas, etc.

**12 Transactions:**
- All Type="sale" from Oct-Nov 2025
- Total amounts range from ₹875 to ₹52,300
- Includes 2 sale_items (Tomato Local 50kg, Onion White 100kg)

**2 Crate Entries:**
- Issue and return transactions

**5 Wastage Logs:**
- Various reasons (expired, damaged, contaminated)

**7 Expiry Alerts:**
- Items expiring in next 1-30 days

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
docker logs pgvms-backend

# Check health
curl http://localhost:8080/api/v1/health

# Restart backend
docker-compose restart backend
```

### "Database connection failed"
```bash
# Check postgres is healthy
docker ps | grep pgvms-postgres

# Check logs
docker logs pgvms-postgres

# Test connection
docker exec pgvms-postgres pg_isready -U pgvms_user -d pgvms
```

### "Frontend shows blank page"
```bash
# Check frontend logs
docker logs pgvms-frontend

# Check if served on port 80 (not 3000 in container)
curl http://localhost:3000

# Rebuild frontend
docker-compose up -d --build frontend
```

### "Charts not showing data"
```bash
# Verify sale_items exist
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM sale_items;"

# Check transaction has Details
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/transactions | jq '.data[0].Details'

# If null, backend needs rebuild
docker-compose up -d --build backend
```

### "No demo data showing"
```bash
# Load demo data
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /docker-entrypoint-initdb.d/demo_simple.sql

# Create demo user
bash setup-demo-user.sh

# Verify data loaded
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM customers;"
```

## File Locations

**Configuration:**
- Environment: `.env` (create from `.env.example`)
- Docker: `docker-compose.yml`
- Nginx: `frontend/nginx.conf`

**Code:**
- Backend entry: `backend/main.go`
- Frontend entry: `frontend/src/App.js`
- API client: `frontend/src/services/api.js`

**Database:**
- Migrations: `infra/migrations/*.sql`
- Demo data: `infra/local/demo_simple.sql`

**Documentation:**
- Setup: `SETUP.md`
- Architecture: `ARCHITECTURE.md`
- API: `docs/API.md`
- This guide: `QUICK_REFERENCE.md`

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

1. ✅ Login with demo@vms.com / demo123
2. ✅ Explore Dashboard - View real charts with sales trends
3. ✅ Browse Inventory - See 45 items with FEFO sorting
4. ✅ View Transactions - Check transaction history with sale items
5. ✅ Create a new customer
6. ✅ Add a new inventory item
7. ✅ Create a sale transaction
8. ✅ Generate a report
9. ✅ Test AI forecasting (if GEMINI_API_KEY set)
10. ✅ Try editing and soft-deleting records

## Support Resources

- **GitHub Issues:** https://github.com/vineethtatipatri19/VMS/issues
- **Full Setup Guide:** [SETUP.md](SETUP.md)
- **Architecture Docs:** [ARCHITECTURE.md](ARCHITECTURE.md)
- **API Reference:** [docs/API.md](docs/API.md)

---

**Quick Start Command:**
```bash
git clone https://github.com/vineethtatipatri19/VMS.git && cd VMS && docker-compose up -d --build && sleep 60 && docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /docker-entrypoint-initdb.d/demo_simple.sql && bash setup-demo-user.sh && echo "✅ Ready! Open http://localhost:3000 and login with demo@vms.com / demo123"
```
