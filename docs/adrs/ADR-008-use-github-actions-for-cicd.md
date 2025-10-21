# ADR-008: Use GitHub Actions for CI/CD

## Status

Accepted

## Date

2025-10-21

## Context

The application requires automated deployment pipelines for:

**Backend**:

- Build Go Lambda functions
- Run tests
- Package Lambda deployment artifacts
- Deploy infrastructure via Terraform
- Deploy Lambda function code

**Frontend**:

- Install npm dependencies
- Build Vue application with Vite
- Run tests and linting
- Deploy static assets to S3
- Invalidate CloudFront cache

CI/CD requirements:

- Automated testing on pull requests
- Automated deployment to production on merge to main
- Separate environments (dev, production)
- Secrets management for AWS credentials
- Cost-effective solution
- Easy to configure and maintain
- Integration with Git workflow

The project will be hosted on GitHub, and we want a streamlined developer experience.

## Decision

We will use **GitHub Actions** for continuous integration and continuous deployment (CI/CD).

Implementation approach:

- Workflow files in `.github/workflows/` directory
- Separate workflows for backend and frontend
- Terraform workflow for infrastructure changes
- AWS credentials stored in GitHub Secrets
- Deployment on push to main branch
- PR checks for testing and validation
- Manual approval gates for production deployments (optional)

## Consequences

### Positive

- **Free for Public Repos**: Unlimited minutes for public repositories
  - **Expected cost: $0/month** if repo is public
  - Private repos: 2,000 minutes/month free (sufficient for low-traffic app)
- **Native GitHub Integration**: Built into GitHub platform
  - No external service setup required
  - Works seamlessly with pull requests
  - Status checks appear directly in PR interface
  - Uses GitHub's authentication and permissions
- **Simple Configuration**: YAML-based workflow files
  - Easy to understand and modify
  - Version-controlled with application code
  - No separate CI/CD platform to learn
- **Rich Ecosystem**: GitHub Actions Marketplace
  - Pre-built actions for common tasks
  - AWS deployment actions maintained by AWS
  - Terraform actions for plan/apply
  - Reduces custom script writing
- **Matrix Builds**: Test against multiple versions/configurations
  - Can test with different Go versions
  - Can test with different Node versions
  - Parallel execution for faster feedback
- **Secrets Management**: Built-in encrypted secrets
  - Store AWS credentials securely
  - Environment-specific secrets
  - No exposure in logs
- **Flexible Triggers**: Run on various Git events
  - Push to branches
  - Pull request creation/updates
  - Manual workflow dispatch
  - Scheduled cron jobs
- **Self-Hosted Runners** (optional): Can use own infrastructure if needed
  - Not needed initially but available if costs grow

### Negative

- **Learning Curve**: Team needs to learn GitHub Actions syntax
  - Workflow YAML structure
  - Action inputs/outputs
  - Context variables and expressions
  - Mitigation: Good documentation; lots of examples available
- **Debugging Difficulty**: Workflows only run in GitHub environment
  - Cannot easily test locally
  - Must push commits to test changes
  - Mitigation: Use `act` tool for local testing; iterate in dev branch
- **Vendor Lock-in**: Tightly coupled to GitHub
  - Migration to GitLab/Bitbucket requires rewriting workflows
  - Mitigation: Low traffic means staying on GitHub is fine
- **Concurrency Limits**: Free tier has limited concurrent jobs
  - Public repos: 20 concurrent jobs
  - Private repos: 5 concurrent jobs (free tier)
  - Mitigation: More than sufficient for this project
- **Execution Time Limits**: 6-hour timeout per job
  - Should never be an issue for our simple builds
  - Mitigation: Builds complete in minutes, not hours
- **No Built-in Artifacts UI**: Artifacts expire after 90 days (default)
  - Must configure retention if needed
  - Mitigation: Deployments go to AWS; don't need long-term artifact storage
- **Cold Start for Runners**: Jobs can take 10-30 seconds to start
  - Slight delay before workflow begins
  - Mitigation: Acceptable for automated deployments

### Cost Analysis

**GitHub Actions Pricing**:

- Public repositories: **FREE** (unlimited minutes)
- Private repositories:
  - Free tier: 2,000 minutes/month
  - Additional: $0.008 per minute (Linux runners)

**Estimated Usage** (private repo scenario):

- Backend build + test + deploy: ~5 minutes
- Frontend build + deploy: ~3 minutes
- Terraform plan/apply: ~5 minutes
- Total per deployment: ~13 minutes

**Deployments per month**: ~20 deployments
**Total minutes**: 20 × 13 = 260 minutes/month

**Expected cost**: **$0/month** (well within 2,000 free minutes)

Even with aggressive development (100 deployments/month):

- Total minutes: 1,300/month
- Cost: $0 (within free tier)

### Workflow Structure

```text
.github/
  workflows/
    backend-ci.yml        # Backend testing on PRs
    backend-deploy.yml    # Backend deployment on merge
    frontend-ci.yml       # Frontend testing on PRs
    frontend-deploy.yml   # Frontend deployment on merge
    terraform-plan.yml    # Infrastructure plan on PRs
    terraform-apply.yml   # Infrastructure apply on merge
```

### Example Workflow: Backend Deployment

```yaml
name: Deploy Backend

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
      - '.github/workflows/backend-deploy.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Run tests
        working-directory: ./backend
        run: go test -v ./...

      - name: Build Lambda
        working-directory: ./backend
        run: |
          GOOS=linux GOARCH=amd64 go build -o bootstrap
          zip lambda.zip bootstrap

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-south-1

      - name: Deploy to Lambda
        run: |
          aws lambda update-function-code \
            --function-name gestione-caselo-graphql \
            --zip-file fileb://backend/lambda.zip
```

### Example Workflow: Frontend Deployment

```yaml
name: Deploy Frontend

on:
  push:
    branches: [main]
    paths:
      - 'frontend/**'
      - '.github/workflows/frontend-deploy.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json

      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci

      - name: Run linting
        working-directory: ./frontend
        run: npm run lint

      - name: Run tests
        working-directory: ./frontend
        run: npm run test

      - name: Build
        working-directory: ./frontend
        run: npm run build

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-south-1

      - name: Deploy to S3
        run: |
          aws s3 sync frontend/dist/ s3://gestione-caselo-frontend/ \
            --delete \
            --cache-control "max-age=31536000,public" \
            --exclude "index.html"

          aws s3 cp frontend/dist/index.html s3://gestione-caselo-frontend/ \
            --cache-control "max-age=0,no-cache,no-store,must-revalidate"

      - name: Invalidate CloudFront
        run: |
          aws cloudfront create-invalidation \
            --distribution-id ${{ secrets.CLOUDFRONT_DISTRIBUTION_ID }} \
            --paths "/*"
```

### Example Workflow: Terraform

```yaml
name: Terraform Apply

on:
  push:
    branches: [main]
    paths:
      - 'terraform/**'

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: 1.6.0

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-south-1

      - name: Terraform Init
        working-directory: ./terraform
        run: terraform init

      - name: Terraform Apply
        working-directory: ./terraform
        run: terraform apply -auto-approve
```

### Security Considerations

**GitHub Secrets** (encrypted environment variables):

- `AWS_ACCESS_KEY_ID`: IAM user access key for deployments
- `AWS_SECRET_ACCESS_KEY`: IAM user secret key
- `CLOUDFRONT_DISTRIBUTION_ID`: CloudFront ID for cache invalidation
- `TERRAFORM_BACKEND_KEY`: Terraform state encryption key (if needed)

**IAM Policy** for deployment user (least-privilege):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "lambda:UpdateFunctionCode",
        "lambda:GetFunction",
        "s3:PutObject",
        "s3:GetObject",
        "s3:ListBucket",
        "s3:DeleteObject",
        "cloudfront:CreateInvalidation",
        "cloudfront:GetInvalidation"
      ],
      "Resource": [
        "arn:aws:lambda:eu-south-1:*:function:gestione-caselo-*",
        "arn:aws:s3:::gestione-caselo-*",
        "arn:aws:s3:::gestione-caselo-*/*",
        "arn:aws:cloudfront::*:distribution/*"
      ]
    }
  ]
}
```

### Workflow Best Practices

1. **Use specific action versions**: `@v4` instead of `@main` (prevents breaking changes)
2. **Cache dependencies**: Use `cache:` for npm, Go modules (faster builds)
3. **Fail fast**: Run tests before deployment steps
4. **Path filters**: Only trigger workflows when relevant files change
5. **Concurrency limits**: Prevent multiple deployments running simultaneously
6. **Environment protection rules**: Require manual approval for production
7. **Status badges**: Add workflow status badges to README
8. **Notification**: Configure Slack/email notifications for failures

## Alternatives Considered

### GitLab CI/CD

- **Pros**:
  - Similar to GitHub Actions
  - Free tier includes 400 minutes/month
  - Good Docker support
  - Can self-host GitLab Runner
- **Cons**:
  - Would need to use GitLab instead of GitHub
  - Team prefers GitHub interface
  - Slightly smaller ecosystem
- **Rejected**: Team preference for GitHub

### CircleCI

- **Pros**:
  - Mature CI/CD platform
  - Good performance
  - Free tier: 6,000 minutes/month
  - Docker-first approach
- **Cons**:
  - External service (another account to manage)
  - More complex configuration
  - Less integration with GitHub
  - Overkill for simple deployment
- **Rejected**: GitHub Actions more integrated; sufficient features

### Jenkins

- **Pros**:
  - Self-hosted (full control)
  - Extremely flexible
  - Large plugin ecosystem
  - Free and open-source
- **Cons**:
  - Must manage Jenkins server (operational burden)
  - Server costs (EC2 instance or always-on machine)
  - Complex setup and maintenance
  - Security updates required
  - Overkill for small project
- **Rejected**: Too much operational overhead; contradicts serverless philosophy

### AWS CodePipeline + CodeBuild

- **Pros**:
  - AWS-native solution
  - Integrates with other AWS services
  - No external dependencies
- **Cons**:
  - Costs $1/active pipeline/month
  - CodeBuild: $0.005/minute (beyond free tier)
  - More expensive than GitHub Actions
  - More complex setup than GitHub Actions
  - Less familiar to team
- **Rejected**: Higher costs; GitHub Actions simpler and free

### Travis CI

- **Pros**:
  - Early CI/CD player
  - Good GitHub integration
- **Cons**:
  - Free tier removed in 2020
  - Now requires paid plan
  - Company had financial troubles
  - Smaller user base now
- **Rejected**: Not free; GitHub Actions better supported

### Bitbucket Pipelines

- **Pros**:
  - Integrated with Bitbucket
  - Free tier: 50 build minutes/month
- **Cons**:
  - Would need to use Bitbucket instead of GitHub
  - Very limited free tier (50 min vs. 2,000 min)
  - Team prefers GitHub
- **Rejected**: Inferior to GitHub Actions

### Manual Deployment Scripts

- **Pros**:
  - Full control
  - No CI/CD platform dependency
  - Simple bash scripts
- **Cons**:
  - Error-prone (human mistakes)
  - No automation or consistency
  - No testing before deployment
  - Slow (manual process)
  - Doesn't scale with team
- **Rejected**: Not sustainable; defeats purpose of CI/CD

## Migration Path

If GitHub Actions becomes insufficient:

1. **Add alternative CI/CD** (e.g., CircleCI, GitLab CI)
2. **Run both in parallel** during transition
3. **Migrate workflow by workflow**
4. **Deprecate GitHub Actions** once migration complete

Most configurations can be adapted with minor YAML changes.

## Monitoring and Alerts

- **Workflow failures**: GitHub sends email notifications
- **Slack integration**: Use GitHub + Slack app for alerts
- **Status badges**: Display build status in README
- **AWS CloudWatch**: Monitor Lambda errors after deployment
- **Manual validation**: Check application after deployment

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)
- [AWS Actions by AWS](https://github.com/aws-actions)
- [Terraform GitHub Actions](https://developer.hashicorp.com/terraform/tutorials/automation/github-actions)
- [GitHub Actions Pricing](https://docs.github.com/en/billing/managing-billing-for-github-actions/about-billing-for-github-actions)
