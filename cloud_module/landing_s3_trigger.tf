resource "aws_lambda_permission" "allow_bucket_trigger_by_landing" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.landing_parse_metadata_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

resource "aws_lambda_permission" "allow_bucket_trigger_by_landing_metadata" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.stories_queue_consumer_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = data.aws_s3_bucket.archive.id

  lambda_function {
    lambda_function_arn = module.landing_parse_metadata_lambda.lambda_function_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "landing.html"
  }

  lambda_function {
    lambda_function_arn = module.stories_queue_consumer_lambda.lambda_function_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "${local.newssite_economy_alias}/"
    filter_suffix       = "/metadata.json"
  }

  depends_on = [
    aws_lambda_permission.allow_bucket_trigger_by_landing,
    aws_lambda_permission.allow_bucket_trigger_by_landing_metadata
  ]
}
