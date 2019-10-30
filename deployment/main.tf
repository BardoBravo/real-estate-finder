variable "bucket" {
  description = "Google Cloud Storage bucket name"
}
variable "project" {
  description = "Google Cloud project ID"
}
variable "region" {
  description = "Google Cloud region"
}
variable "cloud_source_repository" {
  description = "Google Cloud Source repository name"
}
variable "zone" {
  description = "Google Cloud zone, part of the provided region"
}
variable "crontab_schedule" {
  description = "Crontab schedule for running scrapers"
}

locals {
  container_tag = "gcr.io/${var.project}/find-new-rent-service:latest"
}

provider "google" {
  credentials = "${file("credentials.json")}"
  project     = "${var.project}"
  region      = "${var.region}"
  zone        = "${var.zone}"
}

resource "google_storage_bucket" "items" {
  name     = "${var.bucket}"
  location = "US"
}

resource "google_app_engine_application" "app" {
  project     = "${var.project}"
  location_id = "${var.region}"
}

resource "google_cloud_scheduler_job" "job" {
  name        = "run-scrapers"
  description = "Trigger scrapers"
  schedule    = "${var.crontab_schedule}"
  time_zone   = "Etc/UTC"

  http_target {
    http_method = "POST"
    uri         = "${data.external.google_cloud_run_service.result.url}/"
  }

  depends_on = ["google_app_engine_application.app"]
}

resource "google_cloudbuild_trigger" "default" {
  trigger_template {
    branch_name = "master"
    repo_name   = "${var.cloud_source_repository}"
  }

  substitutions = {
    _BUCKET = "${var.bucket}"
  }

  build {
    images = ["${local.container_tag}"]
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "${local.container_tag}", "."]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["push", "${local.container_tag}"]
    }
    step {
      name = "gcr.io/cloud-builders/gcloud"
      args = [
        "beta", "run", "deploy", "find-new-rent-service",
        "--region", "${var.region}",
        "--image", "${local.container_tag}",
        "--update-env-vars", "GCLOUD_BUCKET=$${_BUCKET}",
        "--memory", "1Gi",
        "--timeout", "10m",
        "--platform", "managed",
        "--allow-unauthenticated",
      ]
    }
  }

  provisioner "local-exec" {
    command = "bash trigger_build.sh"
  }
}

data "external" "google_cloud_run_service" {
  depends_on = ["google_cloudbuild_trigger.default"]

  program = ["bash", "get_service_url.sh"]

  query = {
    project = "${var.project}"
    region  = "${var.region}"
  }
}
