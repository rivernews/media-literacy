variable "project_alias" {
  type        = string
  description = "Name prefix used for step function and related resources, including the domain name, so please only use [0-9a-z_-]"
}


variable "slack_signing_secret" {
  type = string
}


variable "slack_post_webhook_url" {
  type = string
}

variable environment_name {
  type = string
  default = ""
  description = "Empty string for Production, otherwise the environment name e.g. dev, stage, etc, make sure to use lowercase (s3 bucket only allows lower)"
}

locals {
  project_name = var.environment_name != "" ? "${var.project_alias}-${var.environment_name}" : "${var.project_alias}"
}
