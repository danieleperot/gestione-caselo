locals {
  budget_monthly = 20
}

# Fetch budget alert email from Parameter Store
data "aws_ssm_parameter" "budget_email" {
  name = "/${var.project}/budget-email"
}

# SNS topic for budget alerts
resource "aws_sns_topic" "budget_alerts" {
  name = "${local.prefix}-budget-alerts"

  tags = {
    Name = "${local.prefix}-budget-alerts"
  }
}

# Email subscription for budget alerts
resource "aws_sns_topic_subscription" "budget_alerts_email" {
  topic_arn = aws_sns_topic.budget_alerts.arn
  protocol  = "email"
  endpoint  = data.aws_ssm_parameter.budget_email.value
}

# Monthly budget with multiple thresholds
resource "aws_budgets_budget" "monthly" {
  name         = "${local.prefix}-monthly-budget"
  budget_type  = "COST"
  limit_amount = local.budget_monthly
  limit_unit   = "USD"
  time_unit    = "MONTHLY"

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 25 # 25% of $20 = $5
    threshold_type             = "PERCENTAGE"
    notification_type          = "ACTUAL"
    subscriber_email_addresses = [data.aws_ssm_parameter.budget_email.value]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 50 # 50% of $20 = $10
    threshold_type             = "PERCENTAGE"
    notification_type          = "ACTUAL"
    subscriber_email_addresses = [data.aws_ssm_parameter.budget_email.value]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 100 # 100% of $20 = $20
    threshold_type             = "PERCENTAGE"
    notification_type          = "ACTUAL"
    subscriber_email_addresses = [data.aws_ssm_parameter.budget_email.value]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 25 # Forecast: 25% of $20 = $5
    threshold_type             = "PERCENTAGE"
    notification_type          = "FORECASTED"
    subscriber_email_addresses = [data.aws_ssm_parameter.budget_email.value]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 100 # Forecast: 100% of $20 = $20
    threshold_type             = "PERCENTAGE"
    notification_type          = "FORECASTED"
    subscriber_email_addresses = [data.aws_ssm_parameter.budget_email.value]
  }

  tags = {
    Name = "${local.prefix}-monthly-budget"
  }
}
