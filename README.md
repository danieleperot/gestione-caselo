# Gestione Caselo

Booking system for Caselo di Salzan, a historic Venetian building used for events.

## Architecture

Built on AWS serverless architecture with:

- Go backend (Lambda + GraphQL)
- Vue 3 frontend (S3 + CloudFront)
- DynamoDB single-table design
- Cognito authentication
- Terraform infrastructure
- GitHub Actions CI/CD

See [Architecture Decision Records](./docs/adrs/README.md) for detailed rationale.

## Local Development

### Prerequisites

- Docker & Docker Compose
- Go 1.25+
- Node.js 20+
- AWS CLI (for LocalStack interaction)
- pre-commit (optional): `pip install pre-commit`

### Setup

```bash
# Install pre-commit hooks (optional but recommended)
pre-commit install

# Start local environment with LocalStack
docker-compose up

# Backend setup
cd backend
go mod download
go test ./...

# Frontend setup
cd frontend
npm install
npm run dev
```

### Development Workflow

See [DEVELOPMENT_RULES.md](./DEVELOPMENT_RULES.md) for detailed development practices.

## Project Structure

```text
gestione-caselo/
├── backend/           # Go Lambda functions
├── frontend/          # Vue 3 application
├── terraform/         # Infrastructure as code
├── docs/
│   └── adrs/         # Architecture Decision Records
├── .github/
│   └── workflows/    # CI/CD pipelines
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
