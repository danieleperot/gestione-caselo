module "shared" {
  source = "../../modules/shared"

  project     = var.project
  environment = var.environment
}

locals {
  prefix = module.shared.prefix
  tags   = module.shared.tags
}
