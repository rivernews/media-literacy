
resource "aws_cloudwatch_event_rule" "scheduler" {
  count = var.environment_name == "" ? 1 : 0

  name                = "${local.project_name}-schedule-start-landings"
  # schedule experssion
  # https://docs.aws.amazon.com/eventbridge/latest/userguide/scheduled-events.html
  schedule_expression = "rate(12 hours)"
  description         = "Starts daily or twice a day so that we get up to date news site changes"
}

resource "aws_cloudwatch_event_target" "scheduler_event_target" {
  count = var.environment_name == "" ? 1 : 0

  target_id = "${local.project_name}-schedule-start-landings-event-target"
  rule      = aws_cloudwatch_event_rule.scheduler.0.name
  arn       = module.pipeline_queue.this_sqs_queue_arn
  sqs_target {
    message_group_id = module.pipeline_queue.this_sqs_queue_name
  }
}

# Scheduler -> SQS
# Based on
# https://github.com/hashicorp/terraform/issues/27347#issuecomment-748961017
# for Scheduler -> Step Function
# see https://stackoverflow.com/questions/65580652/how-to-use-terraform-to-define-cloundwatch-event-rules-to-trigger-stepfunction-s
resource "aws_sqs_queue_policy" "scheduler" {
  count = var.environment_name == "" ? 1 : 0

  queue_url = module.pipeline_queue.this_sqs_queue_id
  policy = data.aws_iam_policy_document.scheduler.0.json
}

data "aws_iam_policy_document" "scheduler" {
  count = var.environment_name == "" ? 1 : 0

  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }

    actions = [
      "sqs:SendMessage",
    ]

    resources = [
      module.pipeline_queue.this_sqs_queue_arn
    ]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values = [
        aws_cloudwatch_event_rule.scheduler.0.arn,
      ]
    }
  }
}


resource "aws_cloudwatch_event_rule" "landing_metadata_scheduler" {
  count = var.environment_name == "" ? 1 : 0

  name                = "${local.project_name}-schedule-start-metadata-for-landing"
  # schedule experssion
  # https://docs.aws.amazon.com/eventbridge/latest/userguide/scheduled-events.html
  schedule_expression = "rate(1 hours)"
  description         = "Every hour to give courtesy to the website"
}

resource "aws_cloudwatch_event_target" "landing_metadata_scheduler_event_target" {
  count = var.environment_name == "" ? 1 : 0

  target_id = "${local.project_name}-schedule-start-metadata-for-landing-event-target"
  rule      = aws_cloudwatch_event_rule.landing_metadata_scheduler.0.name
  arn       = module.landing_metadata_cronjob_lambda.lambda_function_arn
}

module landing_metadata_cronjob_lambda {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-landing-metadata-cronjob-lambda"
  description   = "Query landing pages in db; compute & archive their metadata"
  handler       = "landing_metadata_cronjob"
  runtime       = "go1.x"

  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/landing_metadata_cronjob", ":zip"]
    patterns = ["landing_metadata_cronjob"]
  }]

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  publish = true

  attach_policy_statements = true
  policy_statements = {
    allow_db_query = {
      effect    = "Allow",
      actions   = [
        "dynamodb:Query",
        "dynamodb:UpdateItem",
      ],
      resources = [
        local.media_table_arn,
        "${local.media_table_arn}/index/metadataIndex"
      ]
    }
    s3_archive_bucket = {
      effect    = "Allow",
      actions   = [
        "s3:GetObject",
        "s3:PutObject",
      ],
      resources = [
        "${data.aws_s3_bucket.archive.arn}/*",
      ]
    }
    # enable getting 404 instead of 403 in case of not found
    # https://stackoverflow.com/a/19808954/9814131
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
    LOG_LEVEL = "DEBUG"
    DEBUG = "true"
    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    DYNAMODB_TABLE_ID = local.media_table_id
  }

  tags = {
    Project = local.project_name
  }
}