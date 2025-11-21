# VMS - Complete Setup Guide

This guide will help you get the Vendor Management System running in under 5 minutes on any platform.

## Prerequisites

### All Platforms (Linux, macOS, Windows)

- **Docker Desktop** installed and running
  - Windows: [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/)
  - macOS: [Docker Desktop for Mac](https://docs.docker.com/desktop/install/mac-install/)
  - Linux: [Docker Engine](https://docs.docker.com/engine/install/)
- **Git** installed
  - Windows: [Git for Windows](https://git-scm.com/download/win) or GitHub Desktop
  - macOS: `brew install git` or Xcode Command Line Tools
  - Linux: `sudo apt install git` or `sudo yum install git`
- **4GB RAM** available
- **Ports 3000, 5432, and 8080** available

### Windows-Specific Requirements

**For Windows 10/11:**
1. **WSL 2** (Windows Subsystem for Linux) - Required by Docker Desktop
   - Enable WSL 2: Open PowerShell as Administrator and run:
     ```powershell
     wsl --install
     ```
   - Restart your computer when prompted
   - Set WSL 2 as default: `wsl --set-default-version 2`

2. **Docker Desktop for Windows**
   - Download from [Docker Hub](https://docs.docker.com/desktop/install/windows-install/)
   - During installation, ensure "Use WSL 2 instead of Hyper-V" is selected
   - Start Docker Desktop and wait for it to finish starting (whale icon in system tray)

3. **Git Bash or PowerShell**
   - Git Bash (recommended): Comes with Git for Windows
   - PowerShell: Built into Windows
   - Command Prompt: Also works, but Git Bash preferred for bash scripts

**Verify Installation (Windows):**
```bash
# Open Git Bash or PowerShell and run:
docker --version          # Should show Docker version
docker-compose --version  # Should show Docker Compose version
git --version            # Should show Git version
wsl --status             # Should show WSL 2 is running
```

## 🚀 Quick Start (Recommended)

### One-Command Setup

The easiest way to get started - works on **Linux, macOS, Windows, and GitHub Codespaces**:

**Linux / macOS / Git Bash (Windows):**
```bash
# Clone the repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Run the unified setup script
bash setup.sh
```

**Windows PowerShell:**
```powershell
# Clone the repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Run the unified setup script
bash setup.sh
# OR if bash is not available:
sh setup.sh
```

**Windows Command Prompt:**
```cmd
REM Clone the repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

REM Run setup using Git Bash (recommended)
"C:\Program Files\Git\bin\bash.exe" setup.sh

REM Alternative: Use Docker Compose directly (see Manual Setup section)
```

**That's it!** The script automatically:
- ✅ Detects your environment (local or Codespaces)
- ✅ Configures correct URLs (localhost or Codespaces URLs)
- ✅ Builds all services (backend, frontend, database)
- ✅ Runs 7 database migrations
- ✅ Loads demo data (15 customers, 45 items, 4 transactions)
- ✅ Creates demo user (demo@vms.com / demo123)
- ✅ Verifies everything is working

### Access Your Application

**Local Development:**
- Frontend: http://localhost:3000
- Backend: http://localhost:8080/api/v1

**GitHub Codespaces:**
- URLs displayed in setup output
- Format: `https://<codespace-name>-3000.app.github.dev`
- Click 'Ports' tab → Globe icon next to port 3000

**Login Credentials:**
- Email: `demo@vms.com`
- Password: `demo123`

### What You Get

After setup completes, you'll have:
- **15 Customers**: Mix of B2B, wholesale, and retail
- **45 Inventory Items**: Vegetables, fruits with expiry tracking
- **4 Transactions**: Sales with different payment types
- **5 Crate Ledger Entries**: Issued and returned crates
- **7 Wastage Logs**: Spoiled, damaged, expired items
- **12 Expiry Alerts**: Critical and warning level alerts
- **Full Dashboard**: Real KPIs and charts

---

## Manual Setup (Alternative)

If you prefer step-by-step control or need to troubleshoot:

### 1. Clone and Configure

```bash
# Clone the repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Create environment file (optional - defaults work for local dev)
cat > .env << EOF
JWT_SECRET=local-dev-secret-change-in-production
GEMINI_API_KEY=your-gemini-api-key-optional
MIGRATE_ON_START=true
REACT_APP_API_URL=http://localhost:8080/api/v1
EOF
```

**For Codespaces**, update the API URL:
```bash
# Use your Codespace name
CODESPACE_NAME="your-codespace-name"
echo "REACT_APP_API_URL=https://${CODESPACE_NAME}-8080.app.github.dev/api/v1" >> .env
```

### 2. Start All Services

```bash
# Start the entire stack
docker-compose up -d --build
```

Wait 30-60 seconds for all services to fully start.

### 3. Run Database Migrations

```bash
# Run all migrations (001 through 007)
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done
```

This creates all tables with soft delete support (migration 007).

### 4. Load Demo Data

```bash
# Load customers and inventory
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_simple.sql

# Load transactions, crates, wastage, and alerts
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_additional_data.sql

# Fix sale items data (required for transactions to work)
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "UPDATE sale_items si SET item_name = i.name || ' - ' || i.variant, unit = 'kg' FROM inventory_items i WHERE si.inventory_lot_id = i.id AND si.item_name IS NULL;"
```

### 5. Create Demo User

```bash
# Create a demo user account
bash setup-demo-user.sh
```

Login: **demo@vms.com** / **demo123**

---

## Environment-Specific Notes

### Local Development (Linux / macOS)

No special configuration needed! Just use:
```bash
bash setup.sh
```

Everything runs on localhost with standard ports.

### Windows Local Development

**Using Git Bash (Recommended):**
```bash
# Open Git Bash terminal
cd /c/Users/YourUsername/Projects/VMS
bash setup.sh
```

**Using PowerShell:**
```powershell
# Open PowerShell
cd C:\Users\YourUsername\Projects\VMS
bash setup.sh
# or
sh setup.sh
```

**Using Docker Desktop Directly:**
1. Open Docker Desktop
2. Navigate to project folder in File Explorer
3. Right-click → "Open Git Bash here"
4. Run: `bash setup.sh`

**Common Windows Issues:**

1. **"bash: command not found"**
   - Solution: Use full path to bash: `"C:\Program Files\Git\bin\bash.exe" setup.sh`
   - Or install Git Bash and add to PATH

2. **Line ending issues (CRLF vs LF)**
   ```bash
   # Convert line endings if script fails
   git config --global core.autocrlf false
   git clone https://github.com/vineethtatipatri19/VMS.git
   ```

3. **Permission denied on .sh files**
   ```bash
   # Make scripts executable
   chmod +x setup.sh setup-demo-user.sh
   ```

4. **Docker daemon not running**
   - Start Docker Desktop from Start Menu
   - Wait for whale icon to appear in system tray
   - Ensure WSL 2 backend is running: `wsl --status`

5. **Port already in use**
   ```powershell
   # Check what's using the port (PowerShell as Admin)
   netstat -ano | findstr :3000
   netstat -ano | findstr :8080
   netstat -ano | findstr :5432
   
   # Stop the process (replace PID with actual process ID)
   taskkill /PID <PID> /F
   ```

6. **Docker volumes permission issues**
   - Ensure Docker Desktop has access to your drive
   - Settings → Resources → File Sharing → Add your project folder

**Windows Path Considerations:**
- Use forward slashes `/` in Git Bash: `cd /c/Users/YourName/VMS`
- Use backslashes `\` in PowerShell/CMD: `cd C:\Users\YourName\VMS`
- Docker commands work the same across all terminals

### GitHub Codespaces

Codespaces uses dynamic proxy URLs that change with each session:

**Automatic (Recommended):**
```bash
bash setup.sh  # Automatically detects and configures Codespaces
```

**Manual Configuration:**
```bash
# Get your Codespace name
echo $CODESPACE_NAME

# Your URLs will be:
# Frontend: https://<codespace-name>-3000.app.github.dev
# Backend:  https://<codespace-name>-8080.app.github.dev

# Rebuild frontend with correct API URL
docker-compose build --build-arg REACT_APP_API_URL=https://${CODESPACE_NAME}-8080.app.github.dev/api/v1 frontend
docker-compose up -d frontend
```

**Port Visibility:**
- Ports are automatically forwarded by Codespaces
- If you see a 401/403 error, the port may need to be public
- VS Code usually handles this automatically

---

## Platform-Specific Setup Examples

### Windows 10/11 Complete Walkthrough

**Step 1: Install Prerequisites**
```powershell
# Open PowerShell as Administrator

# Install WSL 2
wsl --install

# Restart computer
# After restart, open PowerShell again

# Download and install Docker Desktop from:
# https://docs.docker.com/desktop/install/windows-install/

# Download and install Git from:
# https://git-scm.com/download/win
```

**Step 2: Setup VMS**
```bash
# Open Git Bash (search in Start Menu)

# Clone repository
cd /c/Users/YourUsername/Documents
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Ensure Docker Desktop is running (check system tray)

# Run setup
bash setup.sh

# Wait 2-3 minutes for all services to start
```

**Step 3: Access Application**
- Open browser: http://localhost:3000
- Login with demo@vms.com / demo123

**Alternative: Use Docker Desktop UI**
1. Open Docker Desktop
2. Click "Images" → "Build" → Select `docker-compose.yml`
3. Click "Containers" to see running services
4. Click port numbers to open in browser

### macOS Complete Walkthrough

**Step 1: Install Prerequisites**
```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Git and Docker
brew install git
brew install --cask docker

# Start Docker Desktop from Applications folder
```

**Step 2: Setup VMS**
```bash
# Open Terminal

# Clone repository
cd ~/Projects
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Run setup
bash setup.sh
```

### Linux (Ubuntu/Debian) Complete Walkthrough

**Step 1: Install Prerequisites**
```bash
# Update system
sudo apt update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt install docker-compose-plugin

# Logout and login again for group changes

# Install Git
sudo apt install git
```

**Step 2: Setup VMS**
```bash
# Clone repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Run setup
bash setup.sh
```

---

## Troubleshooting

### Windows-Specific Issues

**Docker Desktop won't start:**
- Check if WSL 2 is enabled: `wsl --status`
- Ensure virtualization is enabled in BIOS
- Try: Settings → Reset to factory defaults

**"docker-compose: command not found":**
```powershell
# Use docker compose (without hyphen) on newer versions
docker compose up -d

# Or reinstall Docker Desktop
```

**Script fails with "Permission denied":**
```bash
# In Git Bash
chmod +x *.sh
bash setup.sh
```

**Cannot connect to Docker daemon:**
- Ensure Docker Desktop is running
- Check Docker Desktop → Settings → Resources → WSL Integration
- Enable integration for your WSL distribution

**Port already in use (Windows):**
```powershell
# Check what's using the port (PowerShell as Admin)
netstat -ano | findstr :3000
netstat -ano | findstr :8080
netstat -ano | findstr :5432

# Stop the process (replace PID with actual process ID)
taskkill /PID <PID> /F
```

### Services Not Starting

```bash
# Check service status
docker-compose ps

# View logs for a specific service
docker logs pgvms-backend --tail 50
docker logs pgvms-frontend --tail 50
docker logs pgvms-postgres --tail 50

# Restart all services
docker-compose restart
```
docker-compose restart
```

### Database Connection Issues

```bash
# Check if PostgreSQL is ready
docker exec pgvms-postgres pg_isready -U pgvms_user -d pgvms

# Connect to database manually
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms
```

### Port Conflicts

If ports are already in use, edit `docker-compose.yml`:

```yaml
services:
  db:
    ports:
      - "5433:5432"  # Change host port
  backend:
    ports:
      - "8081:8080"  # Change host port
  frontend:
    ports:
      - "3001:80"    # Change host port
```

### Frontend Not Loading

```bash
# Rebuild frontend
docker-compose up -d --build frontend

# Check frontend logs
docker logs -f pgvms-frontend
```

### Backend API Errors

```bash
# Check backend logs
docker logs -f pgvms-backend

# Verify database migrations ran
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"

# Should see these tables:
# customers, inventory_items, transactions, sale_items
# crate_ledger, wastage_log, expiry_alerts
# users, payment_schedules, pricing_tiers, price_history
```

### Migrations Not Running Automatically

**Known Issue:** The backend's automatic migration system (golang-migrate) currently has a path resolution issue.

### Frontend Not Loading / API Errors

**For Local Development:**
```bash
# Verify frontend environment
docker exec pgvms-frontend cat /usr/share/nginx/html/index.html | grep -o 'http://localhost:8080'

# Should return: http://localhost:8080
# If not, rebuild with correct API URL
docker-compose build --build-arg REACT_APP_API_URL=http://localhost:8080/api/v1 frontend
docker-compose up -d frontend
```

**For Codespaces:**
```bash
# Check current API URL
docker logs pgvms-frontend 2>&1 | head -20

# Rebuild with correct Codespaces URL
docker-compose build --no-cache --build-arg REACT_APP_API_URL=https://${CODESPACE_NAME}-8080.app.github.dev/api/v1 frontend
docker-compose up -d frontend
```

### Transactions/Crates/Wastage Not Showing

If you see empty data for these sections:

```bash
# Load additional demo data
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_additional_data.sql

# Fix sale_items (required for transactions)
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "UPDATE sale_items si SET item_name = i.name || ' - ' || i.variant, unit = 'kg' FROM inventory_items i WHERE si.inventory_lot_id = i.id WHERE si.item_name IS NULL;"
```

### Automatic Migrations Not Working

**Known Issue:**

**Symptoms:**
- Backend logs show: `Migration warning: first .: file does not exist`
- Tables are not created automatically
- API returns `column "deleted_at" does not exist` errors
- Database initialization fails when demo SQL files try to insert data before tables exist

**Root Cause:**
The golang-migrate library in the backend expects migration files with `.up.sql` and `.down.sql` suffixes, but our migration files use just `.sql`. Additionally, when demo data SQL files are mounted to the PostgreSQL init directory, they run before tables are created, causing failures.

**Solution 1: Use the Updated setup.sh Script (Recommended)**

The `setup.sh` script now handles migrations correctly:

```bash
bash setup.sh
```

This script:
1. Starts all services WITHOUT demo data auto-loading
2. Waits for database to be ready
3. Runs all 7 migrations manually in order
4. Then loads demo data after tables exist
5. Creates demo user account

**Solution 2: Manual Migration (If setup.sh fails)**

```bash
# 1. Start services
docker-compose up -d

# 2. Wait for database to be ready
sleep 10

# 3. Run all 7 migrations in order
for file in infra/migrations/*.sql; do 
  echo "Applying $(basename $file)..."
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# 4. Verify migrations succeeded
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"

# Should show these tables:
# customers, inventory_items, transactions, sale_items, crate_ledger
# wastage_log, expiry_alerts, users, payment_schedules, pricing_tiers, price_history

# 5. Load demo data
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_simple.sql
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/local/demo_additional_data.sql

# 6. Create demo user
bash setup-demo-user.sh
```

**What Changed in docker-compose.yml:**

The database service no longer auto-mounts demo data during initialization:

```yaml
# OLD (causes issues):
volumes:
  - postgres_data:/var/lib/postgresql/data
  - ./infra/local:/docker-entrypoint-initdb.d:ro  # ❌ Removed

# NEW (correct):
volumes:
  - postgres_data:/var/lib/postgresql/data  # ✅ Only data persistence
```

This prevents demo data SQL from running before migrations create the necessary tables.

**Note:** Migration 007 (`007_add_soft_delete.sql`) is critical - it adds the `deleted_at`, `deleted_by`, and `deletion_reason` columns required by the audit trail system.

### Port Already in Use

```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8080
lsof -i :5432

# Stop the processes or change ports in docker-compose.yml
```

### API Returns Empty Data

If the API returns `{"success": false, "count": 0}`:

1. **Check if demo data was loaded:**
   ```bash
   docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM customers;"
   ```

2. **Verify JWT token is valid:**
   ```bash
   TOKEN=$(cat .demo-token)
   curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers
   ```

3. **Check for soft delete column errors in logs:**
   ```bash
   docker logs pgvms-backend | grep "deleted_at"
   ```

If you see "column deleted_at does not exist", run migration 007 manually.

## Data Reset

To start fresh with clean data:

```bash
# Stop all services and remove volumes
docker-compose down -v

# Start services again
docker-compose up -d --build

# Wait 30 seconds for services to start
sleep 30

# Run migrations manually
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Load demo data
docker cp infra/local/demo_simple.sql pgvms-postgres:/tmp/
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /tmp/demo_simple.sql

# Create demo user
bash setup-demo-user.sh
```

## Manual Setup (Without Docker)

### Backend Setup

1. **Install Go 1.21+**
   ```bash
   # Verify installation
   go version
   ```

2. **Install PostgreSQL 15+**
   ```bash
   # Create database
   createdb pgvms
   ```

3. **Configure Environment**
   ```bash
   export DATABASE_URL="postgres://pgvms_user:pgvms_password@localhost:5432/pgvms?sslmode=disable"
   export JWT_SECRET="your-secret-key"
   export MIGRATE_ON_START="true"
   export PORT="8080"
   ```

4. **Run Backend**
   ```bash
   cd backend
   go mod download
   go run .
   ```

### Frontend Setup

1. **Install Node.js 18+**
   ```bash
   # Verify installation
   node --version
   npm --version
   ```

2. **Install Dependencies**
   ```bash
   cd frontend
   npm install
   ```

3. **Configure Environment**
   ```bash
   # Create .env file
   echo "REACT_APP_API_URL=http://localhost:8080/api/v1" > .env
   ```

4. **Run Frontend**
   ```bash
   npm start
   ```

### Load Demo Data

```bash
# Run SQL file
psql -U pgvms_user -d pgvms -f infra/local/demo_simple.sql

# Create demo user via API
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@vms.com",
    "password": "demo123",
    "name": "Demo User"
  }'
```

## Development Workflow

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Run specific test
go test -v -run TestCustomerCRUD

# Integration tests
cd backend/tests/integration
go test -v
```

### Building for Production

```bash
# Backend
cd backend
go build -o pgvms .

# Frontend
cd frontend
npm run build
```

### Database Migrations

Migrations run automatically when `MIGRATE_ON_START=true`. Manual control:

```bash
# Apply migrations
make migrate-up

# Rollback one migration
make migrate-down

# Check migration status
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT * FROM schema_migrations;"
```

## Useful Commands

### Docker Management

```bash
# View all containers and their status
docker-compose ps

# View real-time logs (all services)
docker-compose logs -f

# View logs for specific service
docker logs pgvms-backend --tail 50 -f
docker logs pgvms-frontend --tail 50 -f
docker logs pgvms-postgres --tail 50 -f

# Stop all services
docker-compose down

# Stop and remove volumes (complete cleanup)
docker-compose down -v

# Restart specific service
docker-compose restart backend

# Rebuild and restart specific service
docker-compose up -d --build backend
```

### Database Access

```bash
# Connect to PostgreSQL CLI
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms

# Quick queries
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM customers;"
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM inventory_items;"
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM transactions;"

# List all tables
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"

# Describe table schema
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\d customers"

# Backup database
docker exec pgvms-postgres pg_dump -U pgvms_user pgvms > backup.sql

# Restore database
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < backup.sql
```

### API Testing

```bash
# Get demo user token
TOKEN=$(cat .demo-token)

# Test various endpoints
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers | jq
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/inventory | jq
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/transactions | jq
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/crates | jq
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/wastage | jq
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/expiry-alerts | jq
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/dashboard | jq

# Check API health
curl http://localhost:8080/api/v1/health
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./... -v

# Run specific test
go test -v -run TestCustomerCRUD

# Integration tests
cd backend/tests/integration
go test -v

# Frontend tests (if configured)
cd frontend
npm test
```

### Development Workflow

```bash
# Make code changes, then:

# Rebuild backend
cd /path/to/VMS
docker-compose up -d --build backend

# Rebuild frontend
docker-compose up -d --build frontend

# Watch logs while developing
docker-compose logs -f backend

# Quick iteration (no cache)
docker-compose build --no-cache backend
docker-compose up -d backend
```

### Database Migrations

```bash
# Run all migrations in order
for file in infra/migrations/*.sql; do 
  echo "Running $(basename $file)..."
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Run specific migration
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/migrations/007_add_soft_delete.sql

# Check migration status (if you have schema_migrations table)
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT * FROM schema_migrations ORDER BY version;"
```

### Building for Production

```bash
# Backend binary
cd backend
go build -o pgvms .
./pgvms

# Frontend production build
cd frontend
npm run build
# Output in build/ directory
```

# Access database
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms

# Check API health
curl http://localhost:8080/api/v1/health

# Test API endpoint (after login)
TOKEN="your-jwt-token"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers
```

## Default Credentials

After running `setup-demo-user.sh`:

**Demo User:**
- Email: `demo@vms.com`
- Password: `demo123`

**Database:**
- Host: `localhost`
- Port: `5432`
- Database: `pgvms`
- User: `pgvms_user`
- Password: `pgvms_password`

## Next Steps

1. **Explore the Dashboard** - View KPIs, sales trends, and top products
2. **Browse Inventory** - See FEFO sorting and expiry alerts
3. **Create a Transaction** - Add a sale with multiple items
4. **Generate Reports** - Try sales, inventory, and customer reports
5. **Test AI Forecasting** - Set GEMINI_API_KEY and predict demand
6. **Customize** - Modify code and see live reload (frontend) or rebuild (backend)

## Production Deployment

See [DEPLOYMENT.md](docs/DEPLOYMENT.md) for:
- Google Cloud Run deployment
- Cloud SQL setup
- CI/CD with Cloud Build
- Environment configuration
- Security best practices

## Support

- **Issues:** https://github.com/vineethtatipatri19/VMS/issues
- **Documentation:** See `/docs` folder
- **API Reference:** [API.md](docs/API.md)

## System Requirements

**Minimum:**
- 2 CPU cores
- 4GB RAM
- 10GB disk space

**Recommended:**
- 4 CPU cores
- 8GB RAM
- 20GB disk space
- SSD storage

## Features Included

- ✅ Complete CRUD operations for all entities
- ✅ Real-time dashboard with charts
- ✅ JWT authentication
- ✅ Soft delete with audit trail
- ✅ Transaction management with sale items
- ✅ Inventory tracking with FEFO sorting
- ✅ Customer management with credit tracking
- ✅ Crate management for returnable assets
- ✅ Wastage logging with photos
- ✅ Expiry alerts with urgency levels
- ✅ AI-powered demand forecasting
- ✅ Comprehensive reporting
- ✅ Print functionality
- ✅ Responsive web design

Happy coding! 🚀
