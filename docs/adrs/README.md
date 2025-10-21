# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records (ADRs) for the Gestione Caselo project.

## What are ADRs?

Architecture Decision Records document important architectural decisions made in the project, including the context, the decision itself, and its consequences. They help team members understand why certain choices were made and provide historical context for future decisions.

## ADR Format

Each ADR follows this structure:

- **Title**: Short descriptive title of the decision
- **Status**: Proposed | Accepted | Deprecated | Superseded
- **Date**: When the decision was made
- **Context**: The issue or problem being addressed
- **Decision**: What was decided and why
- **Consequences**: Positive and negative outcomes of the decision
- **Alternatives Considered**: Other options that were evaluated

## Index of ADRs

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [ADR-001](./ADR-001-use-aws-serverless-architecture.md) | Use AWS Serverless Architecture | Accepted | 2025-10-21 |
| [ADR-002](./ADR-002-use-go-for-backend.md) | Use Go for Backend Lambda Functions | Accepted | 2025-10-21 |
| [ADR-003](./ADR-003-use-graphql-api.md) | Use GraphQL API | Accepted | 2025-10-21 |
| [ADR-004](./ADR-004-use-dynamodb-single-table-design.md) | Use DynamoDB Single-Table Design | Accepted | 2025-10-21 |
| [ADR-005](./ADR-005-use-aws-cognito-for-authentication.md) | Use AWS Cognito for Authentication | Accepted | 2025-10-21 |
| [ADR-006](./ADR-006-use-vue3-with-static-site-generation.md) | Use Vue 3 with Static Site Generation | Accepted | 2025-10-21 |
| [ADR-007](./ADR-007-use-terraform-for-infrastructure.md) | Use Terraform for Infrastructure as Code | Accepted | 2025-10-21 |
| [ADR-008](./ADR-008-use-github-actions-for-cicd.md) | Use GitHub Actions for CI/CD | Accepted | 2025-10-21 |

## Creating a New ADR

1. Copy an existing ADR as a template
2. Number it sequentially (ADR-009, ADR-010, etc.)
3. Use the naming convention: `ADR-XXX-short-title-with-hyphens.md`
4. Fill in all sections with relevant information
5. Update this README with the new ADR in the index table
6. Set status to "Proposed" initially, then update to "Accepted" once approved

## Superseding an ADR

When a decision is reversed or changed:

1. Create a new ADR documenting the new decision
2. Update the old ADR's status to "Superseded by ADR-XXX"
3. Link between the old and new ADRs

## About Gestione Caselo

Gestione Caselo is a booking system for the Caselo di Salzan, a historic Venetian building used for events. The system allows users to book time slots (half-day or full-day) and administrators to manage bookings and schedule events directly.
