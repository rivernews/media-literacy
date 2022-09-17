module "stories_queue" {
  source  = "terraform-aws-modules/sqs/aws"
  version = ">= 2.0, < 3.0"

  # SQS queue attributes: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html

  # FIFO queue should append suffix .fifo
  name = "${local.project_name}-stories-queue"

  delay_seconds = 0

  # so we can use per-message delay
  fifo_queue = false

  # FIFO queue only
  # content_based_deduplication = true

  visibility_timeout_seconds = 3600

  # enable long polling
  receive_wait_time_seconds = 10

  tags = {
    Project = local.project_name
  }
}

module "stories_queue_consumer_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${local.project_name}-fetch-stories"
  description   = "Fetch ${local.project_name} stories; triggered by metadata.json creation"
  handler       = "stories"
  runtime     = "go1.x"
  source_path = [{
    path = "${path.module}/../lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/stories", ":zip"]
    patterns = ["stories"]
  }]
  publish = true

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  reserved_concurrent_executions = -1

  # allow lambda to invoke step function
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

  # event source mapping for long polling
  event_source_mapping = {
    sqs = {
      event_source_arn = module.stories_queue.this_sqs_queue_arn
      batch_size = 1
    }
  }
  allowed_triggers = {
    sqs = {
      principal  = "sqs.amazonaws.com"
      source_arn = module.stories_queue.this_sqs_queue_arn
    }
  }
  attach_policy_statements = true
  policy_statements = {
    pull_sqs = {
      effect    = "Allow",
      actions   = ["sqs:ReceiveMessage", "sqs:DeleteMessage", "sqs:GetQueueAttributes"],
      resources = [module.stories_queue.this_sqs_queue_arn]
    }
    s3_archive_bucket = {
      effect    = "Allow",
      actions   = [
        "s3:PutObject",
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
    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    LOGLEVEL = "DEBUG"
    ENV = local.environment

    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    SFN_ARN = module.batch_stories_sfn.state_machine_arn
  }

  tags = {
    Project = local.project_name
  }
}
