# ADR-001: Use AWS Serverless Architecture

## Status

Accepted

## Date

2025-10-21

## Context

Gestione Caselo is a booking system for the Caselo di Salzan event venue. The application will have very low traffic (few requests per day) as it serves a small local community. The primary requirements are:

- Minimal infrastructure costs
- Reliable availability for booking requests
- Ability to handle future feature additions (pricing, notifications, smart lock integration)
- Simple deployment and maintenance

Traditional server-based hosting (EC2 instances, containers) would incur fixed monthly costs regardless of actual usage, making it inefficient for low-traffic applications.

## Decision

We will use AWS serverless architecture with the following components:

- **AWS Lambda**: Execute backend functions on-demand with pay-per-request pricing
- **Amazon S3**: Host static frontend assets
- **Amazon CloudFront**: CDN for fast global content delivery and HTTPS
- **Amazon DynamoDB**: Serverless NoSQL database with pay-per-request billing
- **Amazon API Gateway**: HTTP/REST API endpoint to trigger Lambda functions
- **Amazon Route 53**: DNS management for custom domain
- **AWS Cognito**: User authentication and authorization

## Consequences

### Positive

- **Cost Efficiency**: With few requests per day, monthly costs will be near-zero (likely under $1/month)
  - Lambda: 1M requests/month free tier, then $0.20 per 1M requests
  - DynamoDB: 25GB storage free, pay-per-request pricing
  - S3: First 5GB free, minimal costs for static assets
  - CloudFront: 1TB data transfer free tier
- **Automatic Scaling**: Handles traffic spikes without configuration (e.g., during popular event announcements)
- **High Availability**: AWS-managed infrastructure with built-in redundancy
- **No Server Maintenance**: No OS patches, security updates, or server management required
- **Fast Development**: Managed services reduce infrastructure complexity

### Negative

- **Cold Starts**: Lambda functions may have 100-500ms latency on first invocation after idle period
  - Mitigation: Acceptable for booking system use case; users won't notice occasional delays
- **Vendor Lock-in**: Deep integration with AWS services makes migration to other platforms difficult
  - Mitigation: Low traffic means costs will remain minimal even long-term
- **Local Development Complexity**: Testing serverless architecture locally requires tools like LocalStack or SAM
- **Lambda Limitations**:
  - 15-minute maximum execution time
  - 10GB memory limit
  - Mitigation: Booking operations are quick; no long-running processes needed
- **Learning Curve**: Team needs to understand serverless patterns and DynamoDB data modeling

### Performance Characteristics

- CloudFront edge caching provides <100ms response times for static assets globally
- Lambda + DynamoDB queries typically complete in 50-200ms (excluding cold starts)
- Cold starts can add 100-500ms for Go Lambda functions

## Alternatives Considered

### Traditional VPS Hosting (DigitalOcean, Linode)

- **Pros**: Simpler mental model, full server control, predictable pricing
- **Cons**: Fixed monthly cost ($5-10/month minimum), requires server maintenance, manual scaling
- **Rejected**: Monthly costs are 5-10x higher for low traffic; server maintenance burden

### Container-based (AWS ECS/Fargate, Google Cloud Run)

- **Pros**: Better for microservices, easier local development, less vendor lock-in
- **Cons**: Minimum monthly cost for running containers (~$15-30/month), more complex than needed
- **Rejected**: Overkill for simple booking system; higher baseline costs

### Platform-as-a-Service (Heroku, Render)

- **Pros**: Very simple deployment, good developer experience
- **Cons**: Minimum $7/month for custom domains, sleep delays on free tier, less control
- **Rejected**: Higher costs than serverless; less flexibility for future integrations (smart locks, etc.)

### Firebase/Supabase

- **Pros**: Backend-as-a-Service with built-in auth, database, and hosting
- **Cons**: Different vendor, less familiarity, potential migration complexity later
- **Rejected**: Team preference for AWS; desire for more control over architecture

## References

- [AWS Lambda Pricing](https://aws.amazon.com/lambda/pricing/)
- [AWS Free Tier](https://aws.amazon.com/free/)
- [Serverless Architecture Patterns](https://aws.amazon.com/serverless/)
