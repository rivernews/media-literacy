resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = data.aws_s3_bucket.archive.id

  lambda_function {
    lambda_function_arn = module.landing_s3_trigger_lambda.lambda_function_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "landing.html"
  }

  lambda_function {
    lambda_function_arn = module.landing_metadata_s3_trigger_lambda.lambda_function_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "/metadata.json"
  }

  depends_on = [
    aws_lambda_permission.allow_bucket_trigger_by_landing,
    aws_lambda_permission.allow_bucket_trigger_by_landing_metadata
  ]
}

resource "aws_lambda_permission" "allow_bucket_trigger_by_landing" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.landing_s3_trigger_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

resource "aws_lambda_permission" "allow_bucket_trigger_by_landing_metadata" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.landing_metadata_s3_trigger_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

module "landing_s3_trigger_lambda" {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-landing-s3-trigger-lambda"
  description   = "Put a landing page in db"
  handler       = "landing_s3_trigger"
  runtime     = "go1.x"

  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/landing_s3_trigger", ":zip"]
    patterns = ["landing_s3_trigger"]
  }]

  timeout = 900
  cloudwatch_logs_retention_in_days = 7
  publish = true

  attach_policy_statements = true
  policy_statements = {
    allow_db_put = {
      effect    = "Allow",
      actions   = [
        "dynamodb:PutItem",
      ],
      resources = [media_table_arn]
    }
  }

  environment_variables = {
    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    LOG_LEVEL = "DEBUG"
    DEBUG = "true"
    DYNAMODB_TABLE_ID = media_table_id
  }

  tags = {
    Project = local.project_name
  }
}

module "landing_metadata_s3_trigger_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${local.project_name}-fetch-stories"
  description   = "Fetch ${local.project_name} stories; triggered by metadata.json creation"
  handler       = "stories"
  runtime     = "go1.x"
  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/stories", ":zip"]
    patterns = ["stories"]
  }]
  publish = true

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  reserved_concurrent_executions = -1

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

  attach_policy_statements = true
  policy_statements = {
    s3_archive_bucket = {
      effect    = "Allow",
      actions   = [
        "s3:GetObject"
      ],
      resources = [
        "${data.aws_s3_bucket.archive.arn}/*",
      ]
    }
    s3_archive_bucket_check_404 = {
      effect    = "Allow",
      actions   = [
        "s3:ListBucket",
      ],
      resources = [
        "${data.aws_s3_bucket.archive.arn}",
      ]
    }
  }

  environment_variables = {
    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    LOGLEVEL = "DEBUG"
    ENV = local.environment

    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    SFN_ARN = module.batch_stories_sfn.state_machine_arn
  }

  tags = {
    Project = local.project_name
  }
}
