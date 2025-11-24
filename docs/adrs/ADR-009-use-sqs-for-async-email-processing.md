# ADR-009: Use SQS for Async Email Processing

## Status

Accepted

## Date

2025-11-24

## Context

The application needs to send notification emails for booking confirmations, approvals, and other events. Sending emails synchronously from the main GraphQL API Lambda would:

- Block user requests while waiting for email delivery
- Make the API vulnerable to email service outages
- Provide no protection against spam/abuse of the email system
- Complicate rate limiting and monitoring

## Decision

We will use **Amazon SQS** (Simple Queue Service) for asynchronous email processing with a dedicated Lambda function.

Architecture:

- EventForm Lambda publishes email messages to SQS queue
- Dedicated Emails Lambda consumes messages from queue
- Lambda configured with **reserved concurrency of 1** (max 1 invocation per second)
- Emails Lambda forwards validated messages to Amazon SES

For local development, we use **ElasticMQ** as an SQS-compatible in-memory queue.

## Consequences

### Positive

- **Non-blocking**: API responds immediately without waiting for email delivery
- **Resilience**: Email service failures don't affect API availability; messages retry automatically
- **Rate Limiting**: Reserved concurrency prevents spam abuse (max 1 email/second)
- **Monitoring**: Queue metrics provide visibility into email processing
  - Queue depth alerts for backlog
  - Dead-letter queue for failed messages
  - CloudWatch metrics for throughput
- **Cost Control**: Rate limiting prevents runaway costs from malicious traffic
- **Separation of Concerns**: Email logic isolated from booking logic

### Negative

- **Eventual Consistency**: Emails arrive after API response (acceptable for notifications)
- **Additional Complexity**: More infrastructure components to manage
- **Local Development**: Requires ElasticMQ container (minimal overhead)

## Alternatives Considered

### Synchronous Email Sending

- **Rejected**: Blocks API, no rate limiting, vulnerable to abuse

### SNS Direct to SES

- **Rejected**: Cannot easily implement per-second rate limiting; less control over retry logic

## References

- [AWS SQS Documentation](https://docs.aws.amazon.com/sqs/)
- [Lambda Reserved Concurrency](https://docs.aws.amazon.com/lambda/latest/dg/configuration-concurrency.html)
- [ElasticMQ](https://github.com/softwaremill/elasticmq)
