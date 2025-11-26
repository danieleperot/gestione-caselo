locals {
  eventform_settings = lookup(var.lambda_settings, "eventform", {
    memory_size = 256
    timeout     = 20
    concurrency = 10
  })
}

# EventForm Lambda (GraphQL API)
resource "aws_lambda_function" "eventform" {
  function_name = "${local.prefix}-eventform"
  role          = aws_iam_role.eventform.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = lookup(local.eventform_settings, "memory_size", 256)
  timeout       = lookup(local.eventform_settings, "timeout", 20)

  architectures = ["x86_64"]
  layers        = []
  ephemeral_storage {
    size = 512
  }

  logging_config {
    log_format            = "JSON"
    log_group             = "/aws/lambda/${local.prefix}-eventform"
    application_log_level = "INFO"
    system_log_level      = "INFO"
  }

  tracing_config {
    mode = "PassThrough"
  }

  # Placeholder S3 location - updated via CI/CD
  s3_bucket = aws_s3_bucket.artifacts.bucket
  s3_key    = "lambda/${local.prefix}-eventform/bootstrap.zip"

  environment {
    variables = {
      DYNAMODB_TABLE = module.dynamodb.main_arn
      SQS_QUEUE_URL  = module.queues.emails_url
      ENVIRONMENT    = var.environment
    }
  }

  tags = {
    Name = "${local.prefix}-eventform"
  }

  lifecycle {
    ignore_changes = [
      s3_key,
      s3_object_version,
      last_modified,
      source_code_hash
    ]
  }
}
