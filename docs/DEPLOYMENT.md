
# Deployment Guide (Cloud Run + Cloud SQL)

1. Build and push container using Cloud Build (cloudbuild.yaml provided) or locally:
   gcloud builds submit --config=infra/cloudbuild.yaml --substitutions=_REGION="us-central1"

2. Create Cloud SQL Postgres instance (or use Terraform in infra/terraform).
3. Set DATABASE_URL to the connection string (use Cloud SQL Proxy or PRIVATE IP + user).
4. Deploy: cloudbuild will deploy to Cloud Run as `pgvms` service.

Notes:
- Ensure the Cloud Run service account has Cloud SQL Client role if connecting via Cloud SQL connector.
- Store JWT secret and DATABASE_URL in Cloud Run environment variables (do NOT hardcode).
