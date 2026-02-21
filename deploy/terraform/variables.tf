variable "project_id" {
    type = string
    description = "GCP Project ID"
    default = "open-telemetry-1000"
}

variable "region" {
    type = string
    description = "GCP Region"
    default = "us-central1"
}


variable "db_tier" {
    type = string
    description = "Cloud SQL Tier"
    default = "db-f1-micro"
}

variable "db_version" {
    type = string
    description = "Cloud SQL Database Version"
    default = "POSTGRES_18"
}

variable "DB_NAME" {
  type        = string
  description = "Cloud SQL database name"
  default     = "otel-mart"
}

variable "DB_USER" {
  type        = string
  description = "Cloud SQL database user"
  default     = "otel-user"
}

variable "DB_PASSWORD" {
  type        = string
  description = "Cloud SQL database password"
  sensitive   = true
}