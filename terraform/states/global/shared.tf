module "shared" {
  source = "../../modules/shared"

  project     = var.project
  environment = "global"
}

locals {
  prefix      = module.shared.prefix
  tags        = module.shared.tags
  github_repo = module.shared.github_repo
}
