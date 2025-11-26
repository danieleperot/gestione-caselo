output "emails_arn" {
  value       = aws_sqs_queue.emails.arn
  description = "ARN of the Emails queue"
}

output "emails_name" {
  value       = aws_sqs_queue.emails.name
  description = "Name of the Emails queue"
}

output "emails_url" {
  value       = aws_sqs_queue.emails.url
  description = "URL of the Emails queue"
}
