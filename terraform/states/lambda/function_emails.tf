locals {
  emails_settings = lookup(var.lambda_settings, "emails", {
    memory_size = 256
    timeout     = 20
    concurrency = 1
  })
}

# Emails Lambda (SQS -> SES)
resource "aws_lambda_function" "emails" {
  function_name = "${local.prefix}-emails"
  role          = aws_iam_role.emails.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = lookup(local.emails_settings, "memory_size", 256)
  timeout       = lookup(local.emails_settings, "timeout", 20)

  architectures = ["x86_64"]
  layers        = []
  ephemeral_storage {
    size = 512
  }

  logging_config {
    log_format            = "JSON"
    log_group             = "/aws/lambda/${local.prefix}-emails"
    application_log_level = "INFO"
    system_log_level      = "INFO"
  }

  tracing_config {
    mode = "PassThrough"
  }

  # Placeholder S3 location - updated via CI/CD
  s3_bucket = aws_s3_bucket.artifacts.bucket
  s3_key    = "lambda/${local.prefix}-emails/bootstrap.zip"

  environment {
    variables = {
      ENVIRONMENT  = var.environment
      ADMIN_EMAILS = "admin@gestione-caselo.it,manager@gestione-caselo.it"
      FROM_ADDRESS = "noreply@gestione-caselo.it"
    }
  }

  tags = {
    Name = "${local.prefix}-emails"
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

# SQS trigger for emails Lambda
resource "aws_lambda_event_source_mapping" "emails_sqs" {
  event_source_arn = module.queues.emails_arn
  function_name    = aws_lambda_function.emails.arn
  batch_size       = 1
  enabled          = true

  scaling_config {
    maximum_concurrency = 2 # Limit concurrent Lambda executions for rate limiting
  }
}
