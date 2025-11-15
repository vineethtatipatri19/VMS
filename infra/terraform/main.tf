
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.0.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region
}

resource "google_project_service" "sqladmin" {
  service = "sqladmin.googleapis.com"
}

resource "google_sql_database_instance" "pgvms_db" {
  name             = "pgvms-db"
  database_version = "POSTGRES_15"
  region           = var.region

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled = false
    }
  }
}

resource "google_sql_database" "pgvms" {
  name     = "pgvms"
  instance = google_sql_database_instance.pgvms_db.name
}

resource "google_service_account" "cloud_run_sa" {
  account_id   = "pgvms-run-sa"
  display_name = "PGVMS Cloud Run Service Account"
}

# NOTE: Deploying Cloud Run with Terraform often requires building a container image in Artifact Registry
# and creating a google_cloud_run_service resource. We provide placeholders and recommend using Cloud Build + gcloud for deployment.
