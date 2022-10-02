module "golang_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${local.project_name}-${local.name}-lambda"
  description   = var.description
  handler       = var.go_handler
  runtime     = "go1.x"
  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build -o ./builds/${var.go_handler} ./cmd/${var.go_handler}", ":zip ./builds/${var.go_handler} ."]
    patterns = ["${var.go_handler}"]
  }]
  publish = true

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  reserved_concurrent_executions = -1

  attach_policy_json = var.attach_policy_json
  policy_json        = var.policy_json

  attach_policy_statements = var.attach_policy_statements
  policy_statements = var.policy_statements

  environment_variables = merge({
    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    ENV = local.environment
    DEBUG = tostring(var.debug)
  }, var.environment_variables)

  tags = {
    Project = local.project_name
  }
}
