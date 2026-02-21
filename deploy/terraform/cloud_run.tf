resource "google_cloud_run_v2_service" "otel-mart-api" {
  name     = "otel-mart-api"
  location = "us-central1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"
  template {
    service_account = google_service_account.otel-mart-sa2.email
    containers {
      name = "otel-mart-api"
      ports {
        container_port = 8080
      }
      image = "us-central1-docker.pkg.dev/${var.project_id}/dilly3-repo/otel-mart-api:latest" 
      volume_mounts {
        name = "empty-dir-volume"
        mount_path = "/mnt"
      }
    }
    volumes {
      name = "empty-dir-volume"
      empty_dir {
        medium = "MEMORY"
        size_limit = "256Mi"
      }
    }
  }
}
