# Based on
# https://github.com/terraform-aws-modules/terraform-aws-lambda/blob/master/examples/build-package/main.tf

module "lambda_layer" {
  source = "terraform-aws-modules/lambda/aws"

  create_layer = true
  layer_name          = "${local.project_name}-lambda-layer"
  description         = "Layer that provides dependencies for the Meida Literacy project"
  runtime     = "python3.8"
  compatible_runtimes = ["python3.8"]
  source_path = [{
    path = "${path.module}/../lambda/layer"
    pip_requirements = true
    # Make sure the follow the Layer Structure
    # https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html#configuration-layers-path
    prefix_in_zip = "python"
  }]

  tags = {
    Project = local.project_name
  }
}


# Based on
# https://github.com/terraform-aws-modules/terraform-aws-step-functions
module "step_function" {
  source = "terraform-aws-modules/step-functions/aws"

  name = "${local.project_name}-step-function"

  # TODO: change to yaml
  definition = templatefile("${path.module}/sfn_def/state_machine_definition.json", {
    SCRAPER_LAMBDA_ARN = module.scraper_lambda.lambda_function_arn
  })

  # allow step function to invoke other service
  #
  # Warning:
  # Needs to create `module.scraper_lambda` before creating this step_function
  # `depends_on` will not help unfortunately
  # https://github.com/terraform-aws-modules/terraform-aws-step-functions/issues/20
  service_integrations = {
    lambda = {
      lambda = [
        module.scraper_lambda.lambda_function_arn
      ]
    }
  }

  type = "STANDARD"

  tags = {
    Project = local.project_name
  }
}

module "scraper_lambda" {
  source = "terraform-aws-modules/lambda/aws"
  create_function = true
  function_name = "${local.project_name}-scraper-lambda"
  description   = "Lambda function for scraping"
  handler       = "landing"
  runtime     = "go1.x"

  # Based on tf https://github.com/terraform-aws-modules/terraform-aws-lambda/blob/master/examples/build-package/main.tf#L111
  # Based on golang https://github.com/snsinfu/terraform-lambda-example/blob/master/Makefile#L23
  source_path = [{
    path = "${path.module}/../lambda_golang/"
    commands = ["go build ./cmd/landing", ":zip"]
    patterns = ["landing"]
  }]

  timeout = 900
  cloudwatch_logs_retention_in_days = 7

  publish = true
  allowed_triggers = {
    # allow sfn to call this func - set from sfn since the sf module provides integration there already
  }

  attach_policy_statements = true
  policy_statements = {
    s3_archive_bucket = {
      effect    = "Allow",
      actions   = [
        "s3:PutObject",
      ],
      resources = ["${data.aws_s3_bucket.archive.arn}/*"]
    }
  }

  environment_variables = {
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
