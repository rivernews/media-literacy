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

module "slack_command_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${local.project_name}-slack-command-lambda"
  description   = "Lambda function for slack command for environment ${local.project_name}"
  handler       = "slack_command_controller.lambda_handler"
  runtime     = "python3.8"
  source_path = "${path.module}/../lambda/src/slack_command_controller.py"

  layers = [
    module.lambda_layer.lambda_layer_arn
  ]

  # Maximum lambda execution time - 15m
  timeout = 20
  cloudwatch_logs_retention_in_days = 7

  # Enable publish to create versions for lambda;
  # otherwise will use $LATEST instead and will cause trouble creating permission for allowing API Gateway invocation:
  # `We currently do not support adding policies for $LATEST.`
  publish = true
  allowed_triggers = {
    APIGatewayAny = {
      service    = "apigateway"
      source_arn = "${module.api.apigatewayv2_api_execution_arn}/*/POST/slack/command"
    }
  }

  attach_policy_statements = true
  policy_statements = {
    pipeline_sqs = {
      effect    = "Allow",
      actions   = ["sqs:SendMessage", "sqs:GetQueueUrl"],
      resources = [module.pipeline_queue.this_sqs_queue_arn]
    }
    s3_archive_bucket = {
      effect    = "Allow",
      actions   = [
        "s3:ListBucket",
      ],
      resources = ["${data.aws_s3_bucket.archive.arn}"]
    }
  }

  # allow lambda to invoke step function
  attach_policy_json = true
  policy_json        = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "states:StartExecution"
            ],
            "Resource": ["${module.batch_stories_sfn.state_machine_arn}"]
        }
    ]
}
EOF

  environment_variables = {
    SLACK_SIGNING_SECRET = var.slack_signing_secret
    SLACK_POST_WEBHOOK_URL = var.slack_post_webhook_url

    PIPELINE_QUEUE_NAME = module.pipeline_queue.this_sqs_queue_name
    BATCH_STORIES_SFN_ARN = module.batch_stories_sfn.state_machine_arn

    LOGLEVEL = "DEBUG"

    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
  }

  tags = {
    Project = local.project_name
  }
}
