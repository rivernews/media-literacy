module main_table {
  source = "../../cloud_module/dynamodb"
  environment_name = var.environment_name
  project_alias = var.project_alias
}