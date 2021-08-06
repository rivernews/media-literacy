data "aws_s3_bucket" "archive" {
  bucket = var.s3_archive_bucket
}
