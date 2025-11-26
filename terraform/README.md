# Terraform Infrastructure

Infrastructure as Code for Gestione Caselo.

## Architecture Overview

```mermaid
flowchart TB
    subgraph Internet
        User[User Browser]
    end

    subgraph Route53["Route53"]
        DNS[DNS: *.gestionecaselo.it]
    end

    subgraph CloudFront["CloudFront Distribution"]
        CF[CDN + SSL Certificate]
    end

    subgraph S3Frontend["S3 Bucket"]
        S3[Frontend Static Files<br/>Vue.js SPA]
    end

    subgraph APIGateway["API Gateway HTTP"]
        APIGW[POST /graphql]
    end

    subgraph Lambda["Lambda Functions"]
        EventForm[EventForm Lambda<br/>GraphQL API<br/>Go Runtime]
        Emails[Emails Lambda<br/>Email Sender<br/>Go Runtime]
    end

    subgraph DynamoDB["DynamoDB"]
        DB[(Single Table<br/>PAY_PER_REQUEST)]
    end

    subgraph SQS["SQS Queues"]
        EmailsQueue[Emails Queue]
        EmailsDLQ[Emails DLQ]
    end

    subgraph SES["Simple Email Service"]
        EmailService[SES Email Delivery]
    end

    subgraph CloudWatch["CloudWatch"]
        Logs[Log Groups]
        Alarms[Alarms & Metrics]
    end

    subgraph Global["Global Resources"]
        OIDC[GitHub OIDC Provider]
        IAM[IAM Roles for CI/CD]
        Budget[Budget Alerts]
        ArtifactsS3[Lambda Artifacts S3]
    end

    User -->|HTTPS| DNS
    DNS -->|Routes to| CF
    CF -->|Serves| S3
    User -->|GraphQL API Calls| APIGW
    APIGW -->|Invokes| EventForm
    EventForm -->|Read/Write| DB
    EventForm -->|Enqueue| EmailsQueue
    EmailsQueue -->|Triggers| Emails
    EmailsQueue -.->|Failed Messages| EmailsDLQ
    Emails -->|Send Email| EmailService
    EventForm -.->|Logs| Logs
    Emails -.->|Logs| Logs
    Logs -.->|Trigger| Alarms
    IAM -.->|Deploy| EventForm
    IAM -.->|Deploy| Emails
    IAM -.->|Upload| S3
    OIDC -.->|Authenticates| IAM

    style User fill:#e1f5ff
    style CF fill:#ff9900
    style S3 fill:#ff9900
    style APIGW fill:#ff9900
    style EventForm fill:#ff9900
    style Emails fill:#ff9900
    style DB fill:#ff9900
    style EmailsQueue fill:#ff9900
    style EmailsDLQ fill:#ff9900
    style EmailService fill:#ff9900
    style DNS fill:#ff9900
    style Logs fill:#ff9900
    style OIDC fill:#ff9900
    style IAM fill:#ff9900
```

## Deployment Strategy

**Infrastructure is deployed manually.** Application code (Lambda, frontend) deploys via GitHub Actions.

## Directory Structure

```text
terraform/
├── tf.sh            # Wrapper script for managing state deployments
├── .env             # Environment configuration (STATE_BUCKET, PROJECT_NAME)
├── .env.example     # Template for .env
├── modules/
│   └── shared/      # Reusable shared module (prefix, tags, etc.)
├── states/          # Per-state Terraform configurations
│   ├── global/      # OIDC provider, IAM roles, budget (one-time, no environment)
│   ├── route53/     # Hosted zones (per environment)
│   ├── cloudfront/  # CloudFront distributions, S3 buckets, certificates (per environment)
│   ├── cloudwatch/  # CloudWatch log groups and alarms (per environment)
│   ├── lambda/      # Lambda functions, IAM roles, SQS queues (per environment)
│   └── api_gateway/ # API Gateway resources (per environment)
└── tfvars/          # Environment-specific variable files
    ├── stage.tfvars
    └── prod.tfvars
```

## Initial Setup

### 1. Create State Bucket and Parameters

```bash
export TF_STATE_BUCKET="YOUR-STATE-BUCKET"

# Create state bucket with versioning, encryption, and block public access
aws s3 mb s3://$TF_STATE_BUCKET --region eu-south-1
aws s3api put-bucket-versioning --bucket $TF_STATE_BUCKET --versioning-configuration Status=Enabled
aws s3api put-bucket-encryption --bucket $TF_STATE_BUCKET \
  --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
aws s3api put-public-access-block --bucket $TF_STATE_BUCKET \
  --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

# Create Parameter Store value for budget alert email
aws ssm put-parameter \
  --name "/gestione-caselo/budget-email" \
  --value "YOUR-EMAIL@example.com" \
  --type "String" \
  --region eu-south-1
```

### 2. Configure Environment

Create `.env` from the example:

```bash
cd terraform
cp .env.example .env
# Edit .env and set:
# STATE_BUCKET=YOUR-STATE-BUCKET
# PROJECT_NAME=gestione-caselo
```

### 3. Deploy Infrastructure

Use `tf.sh` wrapper to deploy states. Deploy in this order:

**Global (one-time, no environment):**

```bash
./tf.sh global init
./tf.sh global plan
./tf.sh global apply
```

**Any other state (per environment):**

```bash
./tf.sh {route53|cloudfront|lambda|cloudwatch|...} {stage|prod} init
./tf.sh {route53|cloudfront|lambda|cloudwatch|...} {stage|prod} plan
./tf.sh {route53|cloudfront|lambda|cloudwatch|...} {stage|prod} apply
```

`./tf.sh` is a wrapper that works with most Terraform commands.

## Naming Convention

Resources use: `${local.prefix}-[resource]` where:

- Global: `prefix = "gestione-caselo-global"`
- Per-env: `prefix = "gestione-caselo-{environment}"`

## Dependencies Between States

Deploy states in this order due to data source dependencies:

1. **global** - No dependencies
2. **route53** - No dependencies (reads global via data source)
3. **cloudfront** - Depends on route53 (hosted zone ID)
4. **cloudwatch** - No dependencies
5. **lambda** - No dependencies (but needs cloudwatch for logs)
6. **api_gateway** - Depends on lambda (function ARNs)
