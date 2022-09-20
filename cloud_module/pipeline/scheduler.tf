
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
