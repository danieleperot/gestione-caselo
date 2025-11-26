variable "environment" {
  description = "Environment name (stage, prod)"
  type        = string
}

variable "project" {
  description = "Project name"
  type        = string
  default     = "gestione-caselo"
}

variable "subdomain" {
  type        = string
  description = "The subdomain from which the application will be served"
  default     = ""
}
