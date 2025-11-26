data "aws_lambda_function" "eventform" {
  function_name = "${local.prefix}-eventform"
}
