locals {
  github_repo = "danieleperot/gestione-caselo"
}

output "prefix" {
  value       = "${var.project}-${var.environment}"
  description = "The string to prefix to the name of all resources created by this module"
}

output "github_repo" {
  value       = local.github_repo
  description = "Path of the Github project"
}

output "tags" {
  value = {
    Environment = var.environment
    Project     = var.project
    ManagedBy   = "terraform"
    Repository  = local.github_repo
  }
  description = "Common tags to apply to all resources"
}

output "domain" {
  value       = "${var.subdomain}.${data.aws_ssm_parameter.base_domain.value}"
  description = "The full domain of the application"
}
