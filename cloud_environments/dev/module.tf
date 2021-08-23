module main {
  source = "../../cloud_module"
  environment_name = var.environment_name
  project_alias = var.project_alias
  slack_signing_secret = var.slack_signing_secret
  slack_post_webhook_url = var.slack_post_webhook_url
}
