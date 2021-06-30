variable "project_name" {
  type        = string
  description = "Name prefix used for step function and related resources, including the domain name, so please only use [0-9a-z_-]"
}


variable "slack_signing_secret" {
  type = string
}
