# ADR-007: Use Terraform for Infrastructure as Code

## Status

Accepted

## Date

2025-10-21

## Context

The application uses multiple AWS services that need to be provisioned and configured:

- Lambda functions (Go runtime)
- API Gateway (GraphQL endpoint)
- DynamoDB table (single-table design)
- Cognito User Pool (authentication)
- S3 buckets (frontend hosting)
- CloudFront distribution (CDN)
- Route 53 (DNS)
- IAM roles and policies

Infrastructure requirements:

- Reproducible environments (dev, staging, production)
- Version-controlled infrastructure changes
- Automated deployment without manual AWS Console clicks
- Easy teardown and recreation for testing
- Cost tracking per environment
- Team collaboration on infrastructure changes

We need an Infrastructure as Code (IaC) solution that is reliable, well-documented, and integrates with CI/CD pipelines.

## Decision

We will use **Terraform** (by HashiCorp) for infrastructure as code.

Implementation approach:

- Terraform configuration files define all AWS resources
- State stored remotely in S3 + DynamoDB for state locking
- Separate workspaces or state files for dev/prod environments
- CI/CD pipeline (GitHub Actions) runs Terraform plan/apply
- Modular structure with reusable modules for common patterns

## Consequences

### Positive

- **Declarative Infrastructure**: Define desired state, Terraform handles the rest
  - Clear, readable configuration files
  - Easy to understand what infrastructure exists
  - Changes are explicit and reviewable
- **Multi-Cloud Support**: Not locked to AWS-specific tools
  - Can add other providers (e.g., Cloudflare, monitoring services)
  - Skills transferable to other cloud providers
  - Future flexibility if architecture changes
- **Excellent AWS Support**: Mature AWS provider with comprehensive resource coverage
  - All AWS services we need are supported
  - Regular updates with new AWS features
  - Large community and extensive documentation
- **State Management**: Tracks resource state and dependencies
  - Knows what exists vs. what should exist
  - Handles resource dependencies automatically
  - Prevents configuration drift
- **Plan Before Apply**: See changes before executing
  - `terraform plan` shows what will change
  - Reduces risk of unexpected modifications
  - Great for code reviews
- **Version Control**: Infrastructure changes tracked in Git
  - Full audit trail of who changed what and when
  - Easy rollback to previous configurations
  - Collaborative infrastructure development
- **Reproducible Environments**: Same code creates identical infrastructure
  - Dev environment mirrors production
  - Easy to spin up test environments
  - Disaster recovery simplified
- **Cost Tracking**: Can tag resources by environment/project
  - Easy to see costs per environment
  - AWS Cost Explorer integration
- **Strong Ecosystem**: Large community and module registry
  - Reusable modules for common patterns
  - Lots of examples and tutorials
  - Third-party tools and integrations

### Negative

- **Learning Curve**: Team needs to learn HCL (Terraform language)
  - Resource syntax and configuration
  - State management concepts
  - Module structure and variables
  - Mitigation: Good documentation; clear examples
- **State File Management**: State file is critical and must be protected
  - Corrupted state can cause serious issues
  - Must use remote state with locking
  - State file contains sensitive data (must be secured)
  - Mitigation: Use S3 backend with versioning + DynamoDB locking
- **Drift Detection**: Manual changes in AWS Console cause drift
  - Terraform state becomes out of sync
  - Requires discipline to always use Terraform
  - Mitigation: Document "no manual changes" policy; use terraform refresh
- **Breaking Changes**: Provider updates can introduce breaking changes
  - Terraform version upgrades may require code changes
  - AWS provider updates occasionally break configurations
  - Mitigation: Pin provider versions; test upgrades in dev environment
- **No Built-in Secrets Management**: Sensitive values need external management
  - Cannot commit passwords/tokens to Git
  - Must use variables or secret management tools
  - Mitigation: Use Terraform variables + AWS Secrets Manager/Parameter Store
- **Apply Time**: Large infrastructure changes can take 10-20 minutes
  - Creating CloudFront distribution: ~10-15 minutes
  - Not instant like manual Console clicks
  - Mitigation: Acceptable trade-off for reproducibility
- **Resource Limits**: Some AWS operations have rate limits
  - Parallel operations may hit throttling
  - Mitigation: Configure parallelism settings

### Project Structure

```text
terraform/
  modules/
    api/              # API Gateway + Lambda module
      main.tf
      variables.tf
      outputs.tf
    frontend/         # S3 + CloudFront module
      main.tf
      variables.tf
      outputs.tf
    database/         # DynamoDB module
      main.tf
      variables.tf
      outputs.tf
  environments/
    dev/
      main.tf
      variables.tf
      terraform.tfvars
      backend.tf
    prod/
      main.tf
      variables.tf
      terraform.tfvars
      backend.tf
  main.tf             # Root module (optional)
  variables.tf
  outputs.tf
  versions.tf         # Terraform and provider versions
```

### State Management

**Remote State Configuration**:

```hcl
terraform {
  backend "s3" {
    bucket         = "gestione-caselo-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "eu-south-1"
    encrypt        = true
    dynamodb_table = "gestione-caselo-terraform-locks"
  }
}
```

**Benefits**:

- State stored in S3 (versioned, encrypted)
- DynamoDB provides state locking (prevents concurrent modifications)
- Team members can collaborate safely
- State backups available via S3 versioning

### Example Resource Definitions

```hcl
# Lambda function
resource "aws_lambda_function" "graphql_api" {
  filename         = "lambda.zip"
  function_name    = "gestione-caselo-graphql"
  role            = aws_iam_role.lambda_role.arn
  handler         = "bootstrap"
  runtime         = "provided.al2"  # Custom runtime for Go
  memory_size     = 256
  timeout         = 30

  environment {
    variables = {
      DYNAMODB_TABLE = aws_dynamodb_table.main.name
      COGNITO_POOL   = aws_cognito_user_pool.main.id
    }
  }

  tags = {
    Environment = "production"
    Project     = "gestione-caselo"
  }
}

# DynamoDB table
resource "aws_dynamodb_table" "main" {
  name           = "gestione-caselo-${var.environment}"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  global_secondary_index {
    name            = "GSI1"
    hash_key        = "GSI1PK"
    range_key       = "GSI1SK"
    projection_type = "ALL"
  }

  tags = {
    Environment = var.environment
    Project     = "gestione-caselo"
  }
}
```

### Workflow

1. **Make Infrastructure Changes**:

   ```bash
   cd terraform/environments/prod
   vim main.tf  # Edit configuration
   ```

2. **Plan Changes**:

   ```bash
   terraform plan -out=tfplan
   # Review what will change
   ```

3. **Apply Changes**:

   ```bash
   terraform apply tfplan
   # Or in CI/CD: terraform apply -auto-approve
   ```

4. **Verify**:

   ```bash
   terraform show  # View current state
   terraform output  # See output values
   ```

### CI/CD Integration

GitHub Actions will:

1. Run `terraform fmt -check` (validate formatting)
2. Run `terraform validate` (validate configuration)
3. Run `terraform plan` (show changes)
4. On merge to main: `terraform apply -auto-approve`

## Alternatives Considered

### AWS CloudFormation

- **Pros**:
  - AWS-native solution
  - No state file management (AWS handles it)
  - ChangeSet preview like Terraform plan
  - Integrated with AWS Console
  - Free (no additional cost)
- **Cons**:
  - JSON/YAML syntax is verbose and harder to read
  - AWS-only (vendor lock-in)
  - Slower to adopt new AWS features than Terraform
  - Less flexible than Terraform HCL
  - Stack limits (500 resources per stack)
  - Circular dependency issues common
- **Rejected**: Less flexible; team preference for Terraform; multi-cloud future possibility

### AWS CDK (Cloud Development Kit)

- **Pros**:
  - Write infrastructure in TypeScript/Python/Go/Java
  - Use familiar programming constructs (loops, conditionals)
  - Better IDE support (autocomplete, refactoring)
  - Synthesizes to CloudFormation
- **Cons**:
  - Still outputs CloudFormation (inherits CFN limitations)
  - Less mature than Terraform
  - Requires knowing a programming language deeply
  - More complex mental model
  - Generated CloudFormation can be hard to debug
  - Team would need to learn CDK on top of Go and Vue
- **Rejected**: Additional complexity; Terraform is simpler and more established

### Pulumi

- **Pros**:
  - Infrastructure as code in real programming languages
  - Strong typing and IDE support
  - Multi-cloud support
  - Modern approach
  - State management similar to Terraform
- **Cons**:
  - Newer than Terraform (less mature)
  - Smaller community
  - Free tier limitations (state storage)
  - Paid service for team features
  - Requires writing code vs. declarative config
- **Rejected**: Terraform more established; free tier more generous

### Serverless Framework

- **Pros**:
  - Simplified Lambda deployment
  - Plugin ecosystem
  - Multi-cloud support
  - Good for function-focused applications
- **Cons**:
  - Limited to serverless resources
  - Still need something else for S3, CloudFront, Cognito, etc.
  - Less comprehensive than Terraform
  - YAML configuration
- **Rejected**: Too limited; would need Terraform anyway for other resources

### AWS SAM (Serverless Application Model)

- **Pros**:
  - AWS-native serverless framework
  - Simplified Lambda and API Gateway definitions
  - Local testing capabilities
  - Free
- **Cons**:
  - Limited to serverless resources (Lambda, API Gateway, DynamoDB)
  - Still need CloudFormation for other resources
  - AWS-only
  - Less flexible than Terraform
- **Rejected**: Too limited; Terraform provides full infrastructure coverage

### Manual AWS Console Configuration

- **Pros**:
  - Immediate visual feedback
  - No learning curve for IaC tools
  - Quick for small changes
- **Cons**:
  - Not reproducible
  - No version control
  - Error-prone (clicking wrong buttons)
  - No audit trail
  - Cannot collaborate effectively
  - Disaster recovery difficult
  - Environment drift inevitable
- **Rejected**: Not sustainable; violates IaC principles

## Migration Path

If Terraform proves insufficient (unlikely):

1. **Export existing resources** with tools like Terraformer
2. **Migrate to alternative IaC** (CloudFormation, CDK, Pulumi)
3. **Test in dev environment** before production
4. **Gradual migration** (can run Terraform + alternative in parallel)

Terraform is well-established, so migration is unlikely to be necessary.

## Best Practices

1. **Always use remote state** with locking
2. **Never commit secrets** to Git
3. **Use consistent naming conventions** for resources
4. **Tag all resources** with Environment, Project, Owner
5. **Pin provider versions** to avoid unexpected changes
6. **Use modules** for reusable patterns
7. **Run terraform plan** before apply
8. **Review plans** in pull requests
9. **Use workspaces or directories** for environments
10. **Document** variable meanings and outputs

## Security Considerations

- **State File Encryption**: Enable S3 encryption at rest
- **State File Access**: Restrict IAM permissions to state bucket
- **Secrets Management**: Use AWS Secrets Manager or SSM Parameter Store
- **IAM Policies**: Follow least-privilege principle
- **Remote State Locking**: Prevent concurrent modifications
- **Audit Logging**: Enable CloudTrail for infrastructure changes

## References

- [Terraform AWS Provider Documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Terraform Best Practices](https://www.terraform-best-practices.com/)
- [Terraform S3 Backend](https://developer.hashicorp.com/terraform/language/settings/backends/s3)
- [AWS Terraform Examples](https://github.com/hashicorp/terraform-provider-aws/tree/main/examples)
