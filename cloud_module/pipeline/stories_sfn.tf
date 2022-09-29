module batch_stories_sfn {
  source = "terraform-aws-modules/step-functions/aws"

  name = "${local.project_name}-batch-stories-sfn"

  definition = templatefile("${path.module}/sfn_def/batch_stories_def.json", {
    FETCH_STORY_LAMBDA_ARN = module.fetch_story_lambda.lambda_function_arn
    STORIES_FINALIZER_LAMBDA_ARN = module.stories_finalizer_lambda.lambda_function_arn
  })

  # allow step function to invoke other service
  #
  # Warning:
  # Needs to create `module.scraper_lambda` before creating this step_function
  service_integrations = {
    lambda = {
      lambda = [
        module.fetch_story_lambda.lambda_function_arn,
        module.stories_finalizer_lambda.lambda_function_arn
      ]
    }
  }

  type = "STANDARD"

  tags = {
    Project = local.project_name
  }
}

module fetch_story_lambda {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-story-lambda"
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
    allow_db_put = {
      effect    = "Allow",
      actions   = [
        "dynamodb:PutItem",
      ],
      resources = [
        local.media_table_arn,
      ]
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

module "stories_finalizer_lambda" {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-stories-finalizer-lambda"
  description   = "Finalizer as last sfn step after all stories fetched"
  handler       = "stories_finalizer"
  runtime     = "go1.x"

  source_path = [{
    path = "${var.repo_dir}/lambda_golang/"
    commands = ["${local.go_build_flags} go build ./cmd/stories_finalizer", ":zip"]
    patterns = ["stories_finalizer"]
  }]

  timeout = 900
  cloudwatch_logs_retention_in_days = 7
  publish = true

  attach_policy_statements = true
  policy_statements = {
    allow_db_put = {
      effect    = "Allow",
      actions   = [
        "dynamodb:Query",
        "dynamodb:UpdateItem",
      ],
      resources = [
        local.media_table_arn,
        "${local.media_table_arn}/index/s3KeyIndex"
      ]
    }
  }

  environment_variables = {
    SLACK_WEBHOOK_URL = var.slack_post_webhook_url
    LOG_LEVEL = "DEBUG"
    DEBUG = "true"
    DYNAMODB_TABLE_ID = local.media_table_id
  }

  tags = {
    Project = local.project_name
  }
}
