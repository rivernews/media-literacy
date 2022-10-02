module batch_stories_sfn {
  source = "terraform-aws-modules/step-functions/aws"

  name = "${local.project_name}-batch-stories-sfn"

  definition = templatefile("${path.module}/sfn_def/batch_stories_def.json", {
    FETCH_STORY_LAMBDA_ARN = "${module.fetch_story_lambda.lambda_function_qualified_arn}"
    STORIES_FINALIZER_LAMBDA_ARN = "${module.stories_finalizer_lambda.lambda_function_qualified_arn}"
  })

  # allow step function to invoke other service
  #
  # Warning:
  # Needs to create `module.scraper_lambda` before creating this step_function
  service_integrations = {
    lambda = {
      lambda = [
        # enforce lambda version; remember to use qualified arn (arn that includes version) for these lambda used in Sfn
        "${module.fetch_story_lambda.lambda_function_arn}:*",
        "${module.stories_finalizer_lambda.lambda_function_arn}:*"
      ]
    }
  }

  type = "STANDARD"

  tags = {
    Project = local.project_name
  }
}

module fetch_story_lambda {
  source = "../media_lambda"
  environment_name = var.environment_name
  project_alias = var.project_alias
  slack_post_webhook_url = var.slack_post_webhook_url
  repo_dir = var.repo_dir

  description = "Fetch and archive a story page"
  go_handler = "story"
  debug = false

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
    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    DYNAMODB_TABLE_ID = local.media_table_id
  }
}

module stories_finalizer_lambda {
  source = "../media_lambda"
  environment_name = var.environment_name
  project_alias = var.project_alias
  slack_post_webhook_url = var.slack_post_webhook_url
  repo_dir = var.repo_dir

  go_handler = "stories_finalizer"
  description = "Finalizer as last sfn step after all stories fetched"
  debug = true

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
    DYNAMODB_TABLE_ID = local.media_table_id
  }
}