terraform {
  backend "s3" {
    # Partial configuration - bucket and key specified at init time
    # terraform init \
    #   -backend-config="bucket=$TF_STATE_BUCKET" \
    #   -backend-config="key=stage/ses.tfstate"
    encrypt      = true
    use_lockfile = true # S3 native locking (Terraform 1.11+)
  }
}
