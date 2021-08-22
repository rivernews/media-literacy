# Based on
# https://github.com/terraform-aws-modules/terraform-aws-apigateway-v2#http-api-gateway
module "api" {
  source = "terraform-aws-modules/apigateway-v2/aws"

  name          = "${local.project_name}-api-gateway"
  description   = "HTTP API Gateway of project ${local.project_name}"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type"]
    allow_methods = ["OPTIONS", "POST", "GET"]
    allow_origins = ["*"]
  }

  # Custom domain
  domain_name                 = local.api_domain_name
  # Note that the certificate has to be in same region if using HTTP API
  domain_name_certificate_arn = aws_acm_certificate_validation.api.certificate_arn

  # Access logs
  default_stage_access_log_destination_arn = aws_cloudwatch_log_group.api.arn
  default_stage_access_log_format          = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage"

  # Routes and integrations
  integrations = {
    "POST /slack/command" = {
      lambda_arn             = module.slack_command_lambda.lambda_function_arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 29000
    }
  }

  default_route_settings = {
    detailed_metrics_enabled = true
    throttling_burst_limit = 5
    throttling_rate_limit = 10
    logging_level = "INFO"
  }

  tags = {
    Project = local.project_name
  }
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "/aws/api/${local.project_name}"
  retention_in_days = 7
}
