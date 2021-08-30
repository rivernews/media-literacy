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

  visibility_timeout_seconds = 60

  # enable long polling
  receive_wait_time_seconds = 10

  tags = {
    Project = local.project_name
  }
}

module "stories_queue_consumer_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${local.project_name}-stories-queue-consumer-lambda"
  description   = "Consumes ${local.project_name} stories queue"
  handler       = "story"
  runtime     = "go1.x"
  source_path = [{
    path = "${path.module}/../lambda_golang/"
    commands = ["go build ./cmd/story", ":zip"]
    patterns = ["story"]
  }]
  publish = true

  timeout = 30
  cloudwatch_logs_retention_in_days = 7

  reserved_concurrent_executions = -1

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
  }

  environment_variables = {
    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    LOGLEVEL = "DEBUG"
    ENV = local.environment
  }

  tags = {
    Project = local.project_name
  }
}
