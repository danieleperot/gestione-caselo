module "dynamodb" {
  source = "./dynamodb"

  prefix = local.prefix
}

module "queues" {
  source = "./queues"

  prefix = local.prefix
}
