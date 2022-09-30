resource "aws_s3_bucket_notification" "bucket_notification_landing_s3_trigger" {
  bucket = data.aws_s3_bucket.archive.id

  lambda_function {
    lambda_function_arn = module.landing_s3_trigger_lambda.lambda_function_arn
    // https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html
    events              = [
      "s3:ObjectCreated:Put",
      "s3:ObjectCreated:CompleteMultipartUpload"
    ]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "/landing.html"
  }

  depends_on = [
    aws_lambda_permission.allow_bucket_trigger_by_landing
  ]
}

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
}
