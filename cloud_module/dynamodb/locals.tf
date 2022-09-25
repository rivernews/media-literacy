locals {
  project_name = var.environment_name != "" ? "${var.project_alias}-${var.environment_name}" : "${var.project_alias}"
  environment = var.environment_name != "" ? var.environment_name : "prod_table"
}
