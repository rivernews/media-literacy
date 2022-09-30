resource "aws_lambda_permission" "allow_bucket_trigger_by_landing_metadata" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.stories_s3_trigger_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

module stories_s3_trigger_lambda {
  source = "../media_lambda"
  environment_name = var.environment_name
  project_alias = var.project_alias
  slack_post_webhook_url = var.slack_post_webhook_url
  repo_dir = var.repo_dir

  description = "Invoke Sfn to fetch all stories; triggered by metadata.json creation"
  go_handler = "stories_s3_trigger"
  debug = false

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
    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    DYNAMODB_TABLE_ID = local.media_table_id
    SFN_ARN = module.batch_stories_sfn.state_machine_arn
  }
}
