data "aws_s3_bucket" "archive" {
  bucket = "${local.project_name}-archives"
}
