output "main_arn" {
  value       = aws_dynamodb_table.main.arn
  description = "ARN of the main table"
}
