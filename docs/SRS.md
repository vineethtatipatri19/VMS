
# PGVMS - Software Requirements Specification (SRS) - Summary

Based on the functional requirements provided by the user. This repo contains:
- Go backend exposing REST endpoints for Inventory, Customers, Transactions, Crates, Forecasting (AI stubs).
- React Native mobile client (skeleton) that consumes the REST API.
- Postgres data model and SQL migrations.
- Dockerfile and Cloud Build configuration to build and deploy to Google Cloud Run.

See `docs/API.md` for detailed API endpoints and `docs/ER.sql` for schema.
