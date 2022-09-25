resource "aws_ssm_parameter" "media_table" {
  name  = "/app/media-literacy/table"
  type  = "String"
  value = "${aws_dynamodb_table.media_table.arn},${aws_dynamodb_table.media_table.id}"
}

// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table#attributes-reference
resource "aws_dynamodb_table" "media_table" {
  name           = "Mediatable"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "uuid"
  range_key      = "createdAt"

  attribute {
    name = "uuid"
    type = "S"
  }

  attribute {
    name = "createdAt"
    type = "S"
  }

  // other fields
  // S3 key
  // docType = {landing | story | landingMetadata | ...}
  // events

  attribute {
    // pontentially sharable by landing AND story pages
    // field `S3 key` will be able to provide newsSiteAlias and landing/story page info
    // value can store doc type, actually
    name = "isDocTypeWaitingForMetadata"
    type = "S"
  }

  // Expiring items by using DynamoDB Time to Live (TTL)
  // https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/TTL.html
  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }

  global_secondary_index {
    // this index is to pull out all landing page that needs metadata generation
    name               = "metadataIndex"
    // must only use equality operator for hash_key
    hash_key           = "isDocTypeWaitingForMetadata"
    // ordering does not matter
    // so range_key need not to be datetime field; (actually its S3 key name already has datetime info)
    // but if we are to specify a sort key field... landing/story distinguish might be good but...
    range_key          = "createdAt"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "INCLUDE"
    non_key_attributes = ["s3Key"]
  }

  tags = {
    Project = local.project_name
    Environment = var.environment_name
  }
}