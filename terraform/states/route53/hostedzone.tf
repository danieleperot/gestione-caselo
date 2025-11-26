resource "aws_route53_zone" "main" {
  name = local.domain

  tags = {
    Name = "${local.prefix}-zone"
  }
}
