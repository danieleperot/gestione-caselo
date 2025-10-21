# TFLint configuration for Terraform validation

plugin "aws" {
  enabled = true
  version = "0.28.0"
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

config {
  module = true
  force  = false
}

rule "terraform_deprecated_interpolation" {
  enabled = true
}

rule "terraform_documented_outputs" {
  enabled = true
}

rule "terraform_documented_variables" {
  enabled = true
}

rule "terraform_naming_convention" {
  enabled = true
}

rule "terraform_required_providers" {
  enabled = true
}

rule "terraform_required_version" {
  enabled = true
}

rule "terraform_unused_declarations" {
  enabled = true
}

rule "terraform_unused_required_providers" {
  enabled = true
}
