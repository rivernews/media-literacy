# only single `aws_s3_bucket_notification` is allowed
resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = data.aws_s3_bucket.archive.id

  lambda_function {
    lambda_function_arn = module.landing_s3_trigger_lambda.lambda_function_arn
    // https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html
    events              = [
      "s3:ObjectCreated:Put",
      "s3:ObjectCreated:CompleteMultipartUpload",
      # "s3:ObjectCreated:*"
    ]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "/landing.html"
  }

  lambda_function {
    lambda_function_arn = module.stories_s3_trigger_lambda.lambda_function_arn
    // https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html
    events              = [
      "s3:ObjectCreated:Put",
      "s3:ObjectCreated:CompleteMultipartUpload",
      # "s3:ObjectCreated:*"
    ]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "/metadata.json"
  }

  depends_on = [
    aws_lambda_permission.allow_bucket_trigger_by_landing,
    aws_lambda_permission.allow_bucket_trigger_by_landing_metadata
  ]
}