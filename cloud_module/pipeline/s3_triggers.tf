resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = data.aws_s3_bucket.archive.id

  lambda_function {
    lambda_function_arn = module.stories_s3_trigger_lambda.lambda_function_arn
    // https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html
    events              = ["s3:ObjectCreated:Put"]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "/metadata.json"
  }

  depends_on = [
    aws_lambda_permission.allow_bucket_trigger_by_landing_metadata
  ]
}

resource "aws_lambda_permission" "allow_bucket_trigger_by_landing_metadata" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.stories_s3_trigger_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

module "stories_s3_trigger_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${local.project_name}-stories-s3-trigger-lambda"
  description   = "Invoke Sfn to fetch all stories; triggered by metadata.json creation"
  handler       = "stories_s3_trigger"
  runtime     = "go1.x"
  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/stories_s3_trigger", ":zip"]
    patterns = ["stories_s3_trigger"]
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
    allow_db_put = {
      effect    = "Allow",
      actions   = [
        "dynamodb:UpdateItem",
      ],
      resources = [
        local.media_table_arn,
      ]
    }
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
    DEBUG = "true"

    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    DYNAMODB_TABLE_ID = local.media_table_id
    SFN_ARN = module.batch_stories_sfn.state_machine_arn
  }

  tags = {
    Project = local.project_name
  }
}
