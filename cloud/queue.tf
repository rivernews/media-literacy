# Based on
# https://github.com/terraform-aws-modules/terraform-aws-sqs/blob/master/examples/complete/main.tf
module "pipeline_queue" {
  source  = "terraform-aws-modules/sqs/aws"
  version = ">= 2.0, < 3.0"

  name = "${var.project_name}-pipeline-queue.fifo"
  receive_wait_time_seconds = 3
  fifo_queue = true
  content_based_deduplication = true
  visibility_timeout_seconds = 30

  tags = {
    Project = var.project_name
  }
}

# Based on
# https://github.com/terraform-aws-modules/terraform-aws-lambda/blob/master/examples/event-source-mapping/main.tf
module "pipeline_queue_consumer_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "${var.project_name}-pipeline-queue-consumer-lambda"
  description   = "Consumer lambda function for ${var.project_name} pipeline queue"
  handler       = "pipeline_queue_consumer.lambda_handler"
  runtime     = "python3.8"
  source_path = "${path.module}/../lambda/src/pipeline_queue_consumer.py"
  publish = true

  layers = [
    module.lambda_layer.lambda_layer_arn
  ]

  cloudwatch_logs_retention_in_days = 7

  # Upstream

  # event source mapping
  event_source_mapping = {
    sqs = {
      event_source_arn = module.pipeline_queue.this_sqs_queue_arn
      # Based on
      # https://github.com/terraform-aws-modules/terraform-aws-lambda/blob/master/main.tf#L222
      batch_size = 1
    }
  }
  allowed_triggers = {
    sqs = {
      principal  = "sqs.amazonaws.com"
      source_arn = module.pipeline_queue.this_sqs_queue_arn
    }
  }
  attach_policy_statements = true
  policy_statements = {
    pull_sqs = {
      effect    = "Allow",
      # Based on
      # https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-api-permissions-reference.html
      actions   = ["sqs:ReceiveMessage", "sqs:DeleteMessage", "sqs:DeleteMessageBatch", "sqs:GetQueueAttributes"],
      resources = [module.pipeline_queue.this_sqs_queue_arn]
    }
  }

  # Downstream

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
            "Resource": ["${module.step_function.state_machine_arn}"]
        }
    ]
}
EOF

  environment_variables = {
    STATE_MACHINE_ARN = module.step_function.state_machine_arn
    SLACK_SIGNING_SECRET = var.slack_signing_secret
    SLACK_POST_WEBHOOK_URL = var.slack_post_webhook_url
    LOGLEVEL = "DEBUG"
  }

  tags = {
    Project = var.project_name
  }
}
