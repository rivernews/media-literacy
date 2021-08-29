module batch_stories_sfn {
  source = "terraform-aws-modules/step-functions/aws"

  name = "${local.project_name}-batch-stories-sfn"

  definition = templatefile("${path.module}/sfn_def/batch_stories_def.json", {
    BATCH_STORIES_FETCH_PARSE_LAMBDA_ARN = module.batch_stories_fetch_parse_lambda.lambda_function_arn
  })

  # allow step function to invoke other service
  #
  # Warning:
  # Needs to create `module.scraper_lambda` before creating this step_function
  service_integrations = {
    lambda = {
      lambda = [
        module.batch_stories_fetch_parse_lambda.lambda_function_arn
      ]
    }
  }

  type = "STANDARD"

  tags = {
    Project = local.project_name
  }
}


module batch_stories_fetch_parse_lambda {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-batch-stories-fetch-parse"
  description   = "Batch fetch and parse all stories of a landing page"
  handler       = "main"
  runtime     = "go1.x"

  source_path = [{
    path = "${path.module}/../scraper_lambda/stories"
    commands = ["go build -o main", ":zip"]
    patterns = ["main"]
  }]

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  publish = true

  attach_policy_statements = true
  policy_statements = {
    pipeline_sqs = {
      effect    = "Allow",
      actions   = ["sqs:SendMessage", "sqs:GetQueueUrl"],
      resources = [module.stories_queue.this_sqs_queue_arn]
    }
    s3_archive_bucket = {
      effect    = "Allow",
      actions   = [
        "s3:PutObject",
      ],
      resources = ["${data.aws_s3_bucket.archive.arn}/*"]
    }
  }

  environment_variables = {
    STORIES_QUEUE_NAME = module.stories_queue.this_sqs_queue_name

    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    LOG_LEVEL = "DEBUG"
    DEBUG = "true"
    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id

    NEWSSITE_ECONOMY = data.aws_ssm_parameter.newssite_economy.value
  }

  tags = {
    Project = local.project_name
  }
}
