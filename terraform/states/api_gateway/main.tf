resource "aws_apigatewayv2_api" "main" {
  name          = "${local.prefix}-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["https://${local.domain}"]
    allow_methods = ["GET", "POST", "OPTIONS"]
    allow_headers = ["content-type", "authorization"]
    max_age       = 300
  }

  tags = {
    Name = "${local.prefix}-api"
  }
}

resource "aws_apigatewayv2_integration" "eventform" {
  api_id           = aws_apigatewayv2_api.main.id
  integration_type = "AWS_PROXY"

  integration_uri        = data.aws_lambda_function.eventform.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "graphql" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "POST /graphql"
  target    = "integrations/${aws_apigatewayv2_integration.eventform.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.main.id
  name        = "$default"
  auto_deploy = true

  default_route_settings {
    throttling_burst_limit = 100
    throttling_rate_limit  = 50
  }

  tags = {
    Name = "${local.prefix}-api-default"
  }
}

# Lambda permission for API Gateway
resource "aws_lambda_permission" "apigw_eventform" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.eventform.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}
