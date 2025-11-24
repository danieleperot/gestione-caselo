# Gestione Caselo

Booking system for Caselo di Salzan, a historic building used for events.

## Architecture

Built on AWS serverless architecture with:

- Go backend (Multiple Lambda functions)
  - EventForm: GraphQL API for bookings
  - Emails: SQS-triggered email processing via SES
- Vue 3 frontend (S3 + CloudFront)
- DynamoDB single-table design
- SQS message queue for async operations
- Cognito authentication
- Terraform infrastructure
- GitHub Actions CI/CD

See [Architecture Decision Records](./docs/adrs/README.md) for detailed rationale.

## Local Development

### Prerequisites

**Required:**

- Docker & Docker Compose
- AWS CLI (for local AWS service container interaction)

**For pre-commit hooks:**

- Python 3 with pip: `pip install pre-commit`
- Go 1.25+
- Node.js 20+ with pnpm
- golangci-lint ([installation](https://golangci-lint.run/welcome/install/))
- Terraform and tflint (if working on infrastructure)

### Setup

```bash
# Install pre-commit hooks (optional but recommended)
pre-commit install

# Start all services (backend, frontend, local AWS service containers)
docker compose up

# Or use watch mode for auto-reload on changes
docker compose watch
```

Access the application:

- Frontend: <http://localhost:5173>
- Backend API (EventForm): <http://localhost:8080>
- ElasticMQ API (SQS): <http://localhost:9324>
- ElasticMQ UI: <http://localhost:9325>
- LocalStack (SES): <http://localhost:4566>
- Mailpit UI (Email viewer): <http://localhost:8025>

### Development Workflow

See [DEVELOPMENT_RULES.md](./DEVELOPMENT_RULES.md) for detailed development practices.

## Project Structure

```text
gestione-caselo/
├── backend/
│   ├── eventform/    # GraphQL API Lambda (bookings)
│   └── emails/       # Email processing Lambda (SQS → SES)
├── frontend/          # Vue 3 application
├── terraform/         # Infrastructure as code
├── docs/
│   └── adrs/         # Architecture Decision Records
├── .github/
│   └── workflows/    # CI/CD pipelines
├── elasticmq.conf    # Local SQS queue definitions
└── docker-compose.yml
```

## Pre-commit Hooks

The project uses pre-commit hooks to ensure code quality:

- **Go**: formatting (gofmt), linting (golangci-lint), tests
- **Frontend**: ESLint, tests
- **Terraform**: formatting, validation, linting
- **General**: trailing whitespace, YAML/JSON validation, secrets detection

## Deployment

Deployments are automated via GitHub Actions on push to `main` branch.

## License

Private project - all rights reserved.
