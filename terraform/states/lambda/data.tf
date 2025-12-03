# Get SES domain identity by domain name
data "aws_ses_domain_identity" "main" {
  domain = module.shared.domain
}
