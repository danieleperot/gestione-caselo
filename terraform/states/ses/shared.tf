module "shared" {
  source = "../../modules/shared"

  project     = var.project
  environment = var.environment
  subdomain   = var.subdomain
}

locals {
  tags   = module.shared.tags
  domain = module.shared.domain
}
