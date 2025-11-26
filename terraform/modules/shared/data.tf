data "aws_ssm_parameter" "base_domain" {
  name = "/${var.project}/base-domain"
}
