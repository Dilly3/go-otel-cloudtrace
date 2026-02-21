resource "google_artifact_registry_repository" "dilly3-repo" {
  location      = "us-central1"
  repository_id = "dilly3-repo"
  description   = "dilly3 docker repository"
  format        = "DOCKER"
}