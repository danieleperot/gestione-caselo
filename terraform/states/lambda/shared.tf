module "shared" {
  source = "../../modules/shared"

  project     = var.project
  environment = var.environment
  subdomain   = var.subdomain
}

locals {
  prefix = module.shared.prefix
  tags   = module.shared.tags
  domain = module.shared.domain
}
