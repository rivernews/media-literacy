module batch_stories_sfn {
  source = "terraform-aws-modules/step-functions/aws"

  name = "${local.project_name}-batch-stories-sfn"

  definition = templatefile("${path.module}/sfn_def/batch_stories_def.json", {
    FETCH_STORY_LAMBDA_ARN = module.fetch_story_lambda.lambda_function_arn
  })

  # allow step function to invoke other service
  #
  # Warning:
  # Needs to create `module.scraper_lambda` before creating this step_function
  service_integrations = {
    lambda = {
      lambda = [
        module.fetch_story_lambda.lambda_function_arn
      ]
    }
  }

  type = "STANDARD"

  tags = {
    Project = local.project_name
  }
}


module landing_parse_metadata_lambda {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-batch-stories-fetch-parse"
  description   = "Scrape metadata from a landing page"
  handler       = "landing_metadata"
  runtime       = "go1.x"

  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/landing_metadata", ":zip"]
    patterns = ["landing_metadata"]
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
        "s3:GetObject"
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

module fetch_story_lambda {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-fetch-story"
  description   = "Fetch and archive a story page"
  handler       = "story"
  runtime       = "go1.x"

  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/story", ":zip"]
    patterns = ["story"]
  }]

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  publish = true

  attach_policy_statements = true
  policy_statements = {
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
  }

  tags = {
    Project = local.project_name
  }
}
