variable "project_alias" {
  type        = string
  description = "Name prefix used for step function and related resources, including the domain name, so please only use [0-9a-z_-]"
}

variable environment_name {
  type = string
  default = ""
  description = "Empty string for Production, otherwise the environment name e.g. dev, stage, etc, make sure to use lowercase (s3 bucket only allows lower)"
}