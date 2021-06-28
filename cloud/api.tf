# HTTP API
resource "aws_apigatewayv2_api" "api" {
	name          = "api-${random_id.id.hex}"
	protocol_type = "HTTP"
	target        = aws_lambda_function.lambda.arn
}

# Permission
resource "aws_lambda_permission" "apigw" {
	action        = "lambda:InvokeFunction"
	function_name = aws_lambda_function.lambda.arn
	principal     = "apigateway.amazonaws.com"

	source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}
