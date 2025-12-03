# Store API Gateway URL in Parameter Store for CI to use
resource "aws_ssm_parameter" "ci_api_gateway_url" {
  name  = "/gestione-caselo/ci/${var.environment}/API_GATEWAY_URL"
  type  = "String"
  value = "${aws_apigatewayv2_api.main.api_endpoint}/graphql"

  tags = {
    Name = "${local.prefix}-ci-api-gateway-url"
  }
}
