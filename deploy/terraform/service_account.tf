
# Enable required APIs
resource "google_project_service" "enabled_services" {
  for_each = toset([
    "compute.googleapis.com",
    "run.googleapis.com",
    "sqladmin.googleapis.com",
    "cloudtrace.googleapis.com",
    "artifactregistry.googleapis.com",
  ])
  project = var.project_id
  service = each.key
}

# Manage existing service account
resource "google_service_account" "otel-mart-sa2" {
  account_id   = "otel-mart-sa2"
  display_name = "Otel Mart Service Account"
}

# Assign IAM roles to the service account
resource "google_project_iam_member" "sa_roles" {
  for_each = toset([
    "roles/compute.admin",            # For Compute Engine
    "roles/run.admin",                # For Cloud Run
    "roles/cloudsql.admin",           # For Cloud SQL
    "roles/cloudtrace.admin",         # For Cloud Trace
    "roles/artifactregistry.admin",   # For Artifact Registry
    "roles/artifactregistry.writer",  # For Artifact Registry
    "roles/artifactregistry.reader",  # For Artifact Registry
    "roles/storage.objectAdmin", # Storage Object
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.otel-mart-sa2.email}"
}
resource "google_project_iam_member" "sa_act_as" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.otel-mart-sa2.email}"
}
resource "google_storage_bucket_iam_member" "terraform_state_access" {
  bucket = "open-telemetry-1000-terraform-state"
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:otel-mart-sa2@open-telemetry-1000.iam.gserviceaccount.com"
}


# Create service account key
resource "google_service_account_key" "otel-mart-sa2" {
  service_account_id = google_service_account.otel-mart-sa2.name
}

# Output the key in base64
output "otel-mart-sa2_key_base64" {
  value     = google_service_account_key.otel-mart-sa2.private_key
  sensitive = true
}

# Write key.json locally
resource "local_file" "sa_key_file2" {
  content  = base64decode(google_service_account_key.otel-mart-sa2.private_key)
  filename = "${path.module}/key2.json"
}

# terraform output -raw otel-mart-sa2_key_base64 | base64 -d > key2.json