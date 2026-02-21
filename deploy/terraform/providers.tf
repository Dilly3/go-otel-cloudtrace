terraform {
  backend "gcs" {
    bucket = "open-telemetry-1000-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = "us-central1"
}