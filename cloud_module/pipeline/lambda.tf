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
    path = "${var.repo_dir}/lambda/layer"
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

module scraper_lambda {
  source = "../media_lambda"
  environment_name = var.environment_name
  project_alias = var.project_alias
  slack_post_webhook_url = var.slack_post_webhook_url
  repo_dir = var.repo_dir

  description = "Fetch landing page"
  go_handler = "landing"
  debug = true

  attach_policy_statements = true
  policy_statements = {
    allow_db_query = {
      effect    = "Allow",
      actions   = [
        "dynamodb:PutItem"
      ],
      resources = [local.media_table_arn]
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
    S3_ARCHIVE_BUCKET = data.aws_s3_bucket.archive.id
    NEWSSITE_ECONOMY = data.aws_ssm_parameter.newssite_economy.value
    DYNAMODB_TABLE_ID = local.media_table_id
  }
}

locals {
  # amd64 is the x86 instruction set
  # arm is not (like M1), not supported by AWS lambda go runtime yet
  # https://stackoverflow.com/questions/26951940/how-do-i-make-go-get-to-build-against-x86-64-instead-of-i386
  go_build_flags = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 "
}
