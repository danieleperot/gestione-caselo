terraform {
  required_version = "~> 1.11"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = local.tags
  }
}

# GitHub OIDC provider for GitHub Actions
# This allows GitHub Actions to assume AWS IAM roles without long-lived credentials
resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = ["sts.amazonaws.com"]

  # GitHub Actions OIDC thumbprints
  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1", # pragma: allowlist secret
    "1c58a3a8518e8759bf075b76b750d4f2df264fcd"  # pragma: allowlist secret
  ]

  tags = {
    Name = "${local.prefix}-github-actions-oidc"
  }
}

# IAM role for Lambda deployment (stage)
resource "aws_iam_role" "github_actions_app_stage" {
  name = "${local.prefix}-github-actions-app-stage"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${local.github_repo}:ref:refs/heads/main"
          }
        }
      }
    ]
  })

  tags = {
    Name = "${local.prefix}-github-actions-app-stage"
  }
}

resource "aws_iam_role_policy" "github_actions_app_stage" {
  name = "${local.prefix}-lambda-update"
  role = aws_iam_role.github_actions_app_stage.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:UpdateFunctionCode",
          "lambda:GetFunction",
          "lambda:PublishVersion"
        ]
        Resource = "arn:aws:lambda:eu-south-1:*:function:gestione-caselo-stage-*"
      }
    ]
  })
}

# IAM role for Lambda deployment (prod)
resource "aws_iam_role" "github_actions_app_prod" {
  name = "${local.prefix}-github-actions-app-prod"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${local.github_repo}:ref:refs/heads/main"
          }
        }
      }
    ]
  })

  tags = {
    Name = "${local.prefix}-github-actions-app-prod"
  }
}

resource "aws_iam_role_policy" "github_actions_app_prod" {
  name = "${local.prefix}-lambda-update"
  role = aws_iam_role.github_actions_app_prod.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:UpdateFunctionCode",
          "lambda:GetFunction",
          "lambda:PublishVersion"
        ]
        Resource = "arn:aws:lambda:eu-south-1:*:function:gestione-caselo-prod-*"
      }
    ]
  })
}

# IAM role for frontend deployment (stage)
resource "aws_iam_role" "github_actions_frontend_stage" {
  name = "${local.prefix}-github-actions-frontend-stage"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${local.github_repo}:ref:refs/heads/main"
          }
        }
      }
    ]
  })

  tags = {
    Name = "${local.prefix}-github-actions-frontend-stage"
  }
}

resource "aws_iam_role_policy" "github_actions_frontend_stage" {
  name = "${local.prefix}-frontend-deploy"
  role = aws_iam_role.github_actions_frontend_stage.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::gestione-caselo-stage-frontend",
          "arn:aws:s3:::gestione-caselo-stage-frontend/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "cloudfront:CreateInvalidation",
          "cloudfront:GetInvalidation",
          "cloudfront:ListDistributions"
        ]
        Resource = "*"
      }
    ]
  })
}

# IAM role for frontend deployment (prod)
resource "aws_iam_role" "github_actions_frontend_prod" {
  name = "${local.prefix}-github-actions-frontend-prod"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${local.github_repo}:ref:refs/heads/main"
          }
        }
      }
    ]
  })

  tags = {
    Name = "${local.prefix}-github-actions-frontend-prod"
  }
}

resource "aws_iam_role_policy" "github_actions_frontend_prod" {
  name = "${local.prefix}-frontend-deploy"
  role = aws_iam_role.github_actions_frontend_prod.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::gestione-caselo-prod-frontend",
          "arn:aws:s3:::gestione-caselo-prod-frontend/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "cloudfront:CreateInvalidation",
          "cloudfront:GetInvalidation",
          "cloudfront:ListDistributions"
        ]
        Resource = "*"
      }
    ]
  })
}
