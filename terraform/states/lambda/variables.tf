variable "project" {
  type        = string
  description = "The name of the project!"
}

variable "environment" {
  type        = string
  description = "The specific environment that we are controlling for the project, like 'stage' or 'prod'"
}

variable "region" {
  type        = string
  description = "AWS region where resources are being deployed"
}

variable "lambda_settings" {
  type = map(object({
    memory_size = number
    timeout     = number
  }))

  description = "Settings for each Lambda function. The key of the map is the Lambda name with no prefixes, and the available settings are `memory_size` for max memory allocated to the Lambda in MB, `timeout` in seconds for max execution time"
}
