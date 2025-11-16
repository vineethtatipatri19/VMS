
# Deployment Guide

## Local Development with Docker Compose (Recommended)

### Prerequisites
- Docker and Docker Compose installed
- Git

### Steps

1. **Clone the repository**
```bash
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS
```

2. **Create environment file**
```bash
cp .env.example .env
```

Edit `.env` and set:
```env
JWT_SECRET=your-secret-key-change-in-production
GEMINI_API_KEY=your-gemini-api-key  # Optional
DATABASE_URL=postgres://pgvms_user:pgvms_password@db:5432/pgvms?sslmode=disable
```

3. **Start all services**
```bash
docker-compose up --build
```

This will start:
- PostgreSQL on port 5432
- Backend API on port 8080
- Frontend on port 3000

4. **Load demo data (optional)**
```bash
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /docker-entrypoint-initdb.d/demo_simple.sql
```

5. **Access the application**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1
- Login: `demo@vms.com` / `demo123`

### Stop services
```bash
docker-compose down
```

### View logs
```bash
docker logs pgvms-backend
docker logs pgvms-frontend
docker logs pgvms-postgres
```

---

## Production Deployment to Google Cloud Run

### Prerequisites
- Google Cloud Project with billing enabled
- `gcloud` CLI installed and authenticated
- Docker installed locally (optional)

### Step 1: Set up Cloud SQL PostgreSQL

1. **Create Cloud SQL instance**
```bash
gcloud sql instances create pgvms-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1 \
  --root-password=STRONG_PASSWORD \
  --database-flags=max_connections=100
```

2. **Create database**
```bash
gcloud sql databases create pgvms --instance=pgvms-db
```

3. **Create user**
```bash
gcloud sql users create pgvms_user \
  --instance=pgvms-db \
  --password=STRONG_PASSWORD
```

4. **Get connection name**
```bash
gcloud sql instances describe pgvms-db --format='value(connectionName)'
# Output: project-id:region:instance-name
```

### Step 2: Build and Push Container Images

1. **Enable required APIs**
```bash
gcloud services enable cloudbuild.googleapis.com run.googleapis.com sqladmin.googleapis.com
```

2. **Build using Cloud Build**
```bash
gcloud builds submit --config=infra/cloudbuild.yaml \
  --substitutions=_REGION="us-central1"
```

Or build locally:
```bash
# Backend
cd backend
docker build -t gcr.io/PROJECT_ID/pgvms-backend .
docker push gcr.io/PROJECT_ID/pgvms-backend

# Frontend
cd ../frontend
docker build -t gcr.io/PROJECT_ID/pgvms-frontend \
  --build-arg REACT_APP_API_URL=https://your-backend-url/api/v1 .
docker push gcr.io/PROJECT_ID/pgvms-frontend
```

### Step 3: Run Database Migrations

1. **Connect to Cloud SQL**
```bash
gcloud sql connect pgvms-db --user=pgvms_user --database=pgvms
```

2. **Run migrations manually or use Cloud SQL proxy**
```bash
# Download Cloud SQL Proxy
curl -o cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.darwin.amd64
chmod +x cloud_sql_proxy

# Start proxy
./cloud_sql_proxy -instances=CONNECTION_NAME=tcp:5432

# In another terminal, run migrations
psql "host=127.0.0.1 port=5432 dbname=pgvms user=pgvms_user" < infra/migrations/001_init.sql
psql "host=127.0.0.1 port=5432 dbname=pgvms user=pgvms_user" < infra/migrations/002_users.sql
psql "host=127.0.0.1 port=5432 dbname=pgvms user=pgvms_user" < infra/migrations/003_add_indexes.sql
psql "host=127.0.0.1 port=5432 dbname=pgvms user=pgvms_user" < infra/migrations/004_enhance_entities.sql
```

### Step 4: Deploy Backend to Cloud Run

```bash
gcloud run deploy pgvms-backend \
  --image gcr.io/PROJECT_ID/pgvms-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --add-cloudsql-instances=CONNECTION_NAME \
  --set-env-vars="DATABASE_URL=postgres://pgvms_user:PASSWORD@/pgvms?host=/cloudsql/CONNECTION_NAME&sslmode=disable" \
  --set-env-vars="JWT_SECRET=your-secret-key" \
  --set-env-vars="PORT=8080" \
  --set-env-vars="GEMINI_API_KEY=your-api-key" \
  --memory=512Mi \
  --cpu=1 \
  --max-instances=10
```

**Get backend URL:**
```bash
gcloud run services describe pgvms-backend --region=us-central1 --format='value(status.url)'
```

### Step 5: Deploy Frontend to Cloud Run

```bash
gcloud run deploy pgvms-frontend \
  --image gcr.io/PROJECT_ID/pgvms-frontend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --memory=256Mi \
  --cpu=1 \
  --max-instances=10
```

### Step 6: Configure Environment

1. **Update backend URL in frontend** (rebuild if needed):
```bash
# Edit frontend Dockerfile or build with:
docker build -t gcr.io/PROJECT_ID/pgvms-frontend \
  --build-arg REACT_APP_API_URL=https://BACKEND_URL/api/v1 \
  frontend/
```

2. **Set up Cloud SQL IAM** (if using Cloud SQL Auth Proxy):
```bash
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:SERVICE_ACCOUNT_EMAIL" \
  --role="roles/cloudsql.client"
```

### Step 7: Load Demo Data (Optional)

```bash
# Using Cloud SQL Proxy
psql "host=127.0.0.1 port=5432 dbname=pgvms user=pgvms_user" < infra/local/demo_simple.sql
```

---

## Environment Variables Reference

### Backend (Cloud Run)
| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Yes | `postgres://user:pass@/db?host=/cloudsql/...` |
| `PORT` | Server port | Yes | `8080` |
| `JWT_SECRET` | Secret for JWT tokens | Yes | `your-secret-key` |
| `GEMINI_API_KEY` | Google Gemini API key | No | `AIza...` |
| `MIGRATE_ON_START` | Auto-run migrations | No | `true` |

### Frontend (Build Args)
| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `REACT_APP_API_URL` | Backend API URL | Yes | `https://backend.run.app/api/v1` |

---

## Security Best Practices

1. **Use Secret Manager for sensitive values**
```bash
# Store secret
gcloud secrets create jwt-secret --data-file=-
# Grant access
gcloud secrets add-iam-policy-binding jwt-secret \
  --member="serviceAccount:SERVICE_ACCOUNT" \
  --role="roles/secretmanager.secretAccessor"
```

2. **Enable Cloud Armor** for DDoS protection

3. **Set up VPC connector** for private Cloud SQL access

4. **Use Cloud Run service accounts** with least privilege

5. **Enable audit logging**

6. **Use HTTPS only** (enforced by default on Cloud Run)

---

## Monitoring & Logging

### View Logs
```bash
# Backend logs
gcloud run logs read --service=pgvms-backend --limit=50

# Frontend logs
gcloud run logs read --service=pgvms-frontend --limit=50
```

### Set up Monitoring
```bash
# Enable Cloud Monitoring
gcloud services enable monitoring.googleapis.com

# Create uptime check
gcloud monitoring uptime create pgvms-uptime \
  --resource-type=uptime-url \
  --http-check-path=/api/v1/health \
  --http-check-method=GET
```

---

## Scaling Configuration

Cloud Run auto-scales based on traffic. Configure:

```bash
gcloud run services update pgvms-backend \
  --min-instances=1 \
  --max-instances=10 \
  --concurrency=80 \
  --cpu-throttling \
  --region=us-central1
```

---

## Backup & Recovery

### Automated Backups (Cloud SQL)
```bash
gcloud sql instances patch pgvms-db \
  --backup-start-time=02:00 \
  --enable-bin-log
```

### Manual Backup
```bash
gcloud sql backups create --instance=pgvms-db
```

### Restore from Backup
```bash
gcloud sql backups list --instance=pgvms-db
gcloud sql backups restore BACKUP_ID --backup-instance=pgvms-db
```

---

## Troubleshooting

### Backend won't start
- Check Cloud SQL connection: verify CONNECTION_NAME
- Check logs: `gcloud run logs read --service=pgvms-backend`
- Verify service account has Cloud SQL Client role

### Frontend can't reach backend
- Verify REACT_APP_API_URL is set correctly
- Check CORS settings in backend
- Verify backend is publicly accessible

### Database connection issues
- Verify Cloud SQL instance is running
- Check VPC connector configuration
- Verify database credentials

---

## Cost Optimization

1. Use **f1-micro** Cloud SQL tier for dev/staging
2. Set **min-instances=0** for non-production environments
3. Enable **request-based pricing** on Cloud Run
4. Use **Cloud Storage** for static assets
5. Implement **caching** with Cloud CDN

---

## Rollback

```bash
# List revisions
gcloud run revisions list --service=pgvms-backend --region=us-central1

# Rollback to previous revision
gcloud run services update-traffic pgvms-backend \
  --to-revisions=REVISION_NAME=100 \
  --region=us-central1
```
