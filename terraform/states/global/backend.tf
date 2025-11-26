terraform {
  backend "s3" {
    # Partial configuration - bucket specified at init time
    # terraform init -backend-config="bucket=$TF_STATE_BUCKET"
    key          = "global/oidc.tfstate"
    region       = "eu-south-1"
    encrypt      = true
    use_lockfile = true # S3 native locking (Terraform 1.11+)
  }
}
