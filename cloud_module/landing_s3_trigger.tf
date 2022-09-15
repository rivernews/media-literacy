resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.batch_stories_fetch_parse_lambda.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.archive.arn
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = data.aws_s3_bucket.archive.id

  lambda_function {
    lambda_function_arn = module.batch_stories_fetch_parse_lambda.lambda_function_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "media-literacy-archives/"
    filter_suffix       = ".html"
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}
