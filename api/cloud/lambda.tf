# Based on
# https://github.com/terraform-aws-modules/terraform-aws-lambda/blob/master/examples/build-package/main.tf

module "slack_command_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  create_function = true
  function_name = "slack-command-lambda"
  description   = "Lambda function for slack command"
  handler       = "slack_command_controller.lambda_handler"
  runtime     = "python3.8"
  source_path = "${path.module}/../lambda/src"

  layers = [
    module.lambda_layer.lambda_layer_arn
  ]

  cloudwatch_logs_retention_in_days = 7

  # Enable publish to create versions for lambda;
  # otherwise will use $LATEST instead and will cause trouble creating permission for allowing API Gateway invocation:
  # `We currently do not support adding policies for $LATEST.`
  publish = true
  allowed_triggers = {
    APIGatewayAny = {
      service    = "apigateway"
      source_arn = "${module.api.apigatewayv2_api_execution_arn}/*/POST/slack/command"
    }
  }

  environment_variables = {
    SLACL_SIGNING_SECRET = var.slack_signing_secret
  }

  tags = {
    Project = var.project_name
  }
}

module "lambda_layer" {
  source = "terraform-aws-modules/lambda/aws"

  create_layer = true
  layer_name          = "${var.project_name}-lambda-layer"
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
    Project = var.project_name
  }
}
