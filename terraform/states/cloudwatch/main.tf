resource "aws_cloudwatch_log_group" "eventform" {
  name              = "/aws/lambda/${local.prefix}-eventform"
  retention_in_days = 30

  tags = {
    Name = "${local.prefix}-eventform-logs"
  }
}

resource "aws_cloudwatch_log_group" "emails" {
  name              = "/aws/lambda/${local.prefix}-emails"
  retention_in_days = 30

  tags = {
    Name = "${local.prefix}-emails-logs"
  }
}
