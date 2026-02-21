# 1. Cloud SQL Instance
resource "google_sql_database_instance" "otel-mart-db-instance" {
  name             = var.DB_NAME
  database_version = var.db_version
  region           = var.region

  settings {
    tier = var.db_tier
    edition = "ENTERPRISE"
  }
}

# 2. Database
resource "google_sql_database" "otel-mart-db" {
  name     = var.DB_NAME
  instance = google_sql_database_instance.otel-mart-db-instance.name
}

# 3. User
resource "google_sql_user" "otel-mart-db-user" {
  name     = var.DB_USER
  instance = google_sql_database_instance.otel-mart-db-instance.name
  password = var.DB_PASSWORD
}

# 5. Allow unauthenticated access
resource "google_cloud_run_v2_service_iam_member" "noauth" {
  name     = google_cloud_run_v2_service.otel-mart-api.name
  location = google_cloud_run_v2_service.otel-mart-api.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
# Output the IP address (host)
output "db_host" {
  value = google_sql_database_instance.otel-mart-db-instance.public_ip_address
}

# Output the connection name (useful for Cloud SQL Proxy / Cloud Run)
output "db_connection_name" {
  value = google_sql_database_instance.otel-mart-db-instance.connection_name
}