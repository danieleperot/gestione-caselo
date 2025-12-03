# Store allowed origins in Parameter Store for Lambda to use
resource "aws_ssm_parameter" "allowed_origins" {
  name  = "/gestione-caselo/${var.environment}/ALLOWED_ORIGINS"
  type  = "String"
  value = "https://${local.domain}"

  tags = {
    Name = "${local.prefix}-allowed-origins"
  }
}
