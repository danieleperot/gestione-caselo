
# Dead letter queue for emails
resource "aws_sqs_queue" "emails_dlq" {
  name                      = "${var.prefix}-emails-dlq"
  message_retention_seconds = 1209600 # 14 days

  tags = {
    Name = "${var.prefix}-emails-dlq"
  }
}

# Main emails queue
resource "aws_sqs_queue" "emails" {
  name                       = "${var.prefix}-emails"
  visibility_timeout_seconds = 300
  message_retention_seconds  = 345600 # 4 days

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.emails_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.prefix}-emails"
  }
}
