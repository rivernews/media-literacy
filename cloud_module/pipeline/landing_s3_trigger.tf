resource "aws_lambda_permission" "allow_bucket_trigger_by_landing" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.landing_s3_trigger_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

module landing_s3_trigger_lambda {
  source = "../media_lambda"
  environment_name = var.environment_name
  project_alias = var.project_alias
  slack_post_webhook_url = var.slack_post_webhook_url
  repo_dir = var.repo_dir

  description = "Triggered by landing creation; put landing to db"
  go_handler = "landing_s3_trigger"
  debug = true

  attach_policy_statements = true
  policy_statements = {
    allow_db_query = {
      effect    = "Allow",
      actions   = [
        "dynamodb:PutItem"
      ],
      resources = [local.media_table_arn]
    }
  }

  environment_variables = {
    DYNAMODB_TABLE_ID = local.media_table_id
  }
}
