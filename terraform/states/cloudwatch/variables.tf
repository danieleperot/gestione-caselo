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
