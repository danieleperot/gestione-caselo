# Development Rules - Gestione Caselo

This document defines development practices and conventions for working with Claude on this project.

## Architecture Context

This project follows the architecture decisions documented in `/docs/adrs/`:

- AWS Serverless (Lambda + DynamoDB + S3 + CloudFront)
- Go backend with GraphQL API
- Vue 3 frontend with Tailwind CSS
- Terraform for infrastructure
- GitHub Actions for CI/CD

Refer to ADRs for detailed rationale and trade-offs.

## Development Workflow

### Local Development with Docker

All local development uses **Docker + Docker Compose** for consistency and isolation:

- **LocalStack Community**: Emulate AWS services locally (DynamoDB, S3, Lambda, Cognito, etc.)
- **Use LocalStack whenever possible** for testing AWS integrations
- **Docker Compose**: Orchestrate all services (LocalStack, backend, frontend, databases)
- **No AWS credentials needed** for local development

Typical docker-compose.yml services:

- `localstack`: AWS service emulation
- `backend`: Go Lambda functions running locally
- `frontend`: Vue dev server with hot reload
- `dynamodb-admin` (optional): UI for viewing DynamoDB tables

### Test-Driven Learning Approach

The preferred development workflow is:

1. **I write unit tests first** - Covering the expected behavior
2. **You implement production code** - To make tests pass
3. **You learn by doing** - Writing code to solve test requirements
4. **I provide guidance** - Via docs, examples, or explanations when needed

### Code-First Philosophy

You prefer to code first and learn through implementation. When you need help:

- Ask for Go documentation links or explanations
- Request examples of patterns you're implementing
- Ask for guidance on testing strategies

I will provide pointers to relevant docs rather than writing extensive code unless explicitly requested.

## Communication Style

### Tone Guidelines

- **Factual and direct** - Focus on technical information
- **Minimal enthusiasm** - Avoid phrases like "Absolutely!", "Great!", "You're doing amazing!"
- **Introvert-friendly** - Concise, no unnecessary cheerleading
- **Objective feedback** - Point out issues directly without softening language excessively

### Good Examples

✓ "This implementation has a race condition on line 45."
✓ "The test fails because the mock isn't configured for that case."
✓ "Consider using `sync.WaitGroup` here."

### Avoid

✗ "Great work! This is looking amazing!"
✗ "You're absolutely crushing it!"
✗ "I'm so excited to help you with this!"

## Technology Stack Conventions

### Backend (Go)

- **Runtime**: Go 1.25+
- **Lambda handler**: Use `github.com/aws/aws-lambda-go/lambda`
- **DynamoDB**: Use AWS SDK for Go v2
- **GraphQL**: Use `gqlgen` for schema-first development
- **Testing**: Standard `testing` package, table-driven tests preferred
- **Error handling**: Explicit error returns, wrap with context

**Key patterns**:

```go
// Table-driven tests
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test body
        })
    }
}
```

### Frontend (Vue 3)

- **Framework**: Vue 3 Composition API (not Options API)
- **Build tool**: Vite
- **Styling**: Tailwind CSS utility classes
- **State**: Pinia for global state
- **GraphQL**: Apollo Client or urql
- **Testing**: Vitest for unit tests, Playwright for e2e

**Key patterns**:

```vue
<script setup lang="ts">
// Composition API - use this approach
import { ref, computed } from 'vue'
</script>
```

### Infrastructure (Terraform)

- **Version**: Pin Terraform and provider versions
- **State**: Remote state in S3 with DynamoDB locking
- **Structure**: Modules for reusable components
- **Naming**: `{project}-{environment}-{resource}` convention

### DynamoDB Schema

- **Single-table design** as per ADR-004
- **Key prefixes**: USER#, BOOKING#, EVENT#, etc.
- **Access patterns**: Document all queries before implementing
- **Indexes**: Use GSI1, GSI2 for alternate access patterns

Example entity structure:

```text
PK: USER#{userId}
SK: METADATA
GSI1PK: EMAIL#{email}
GSI1SK: USER
```

## Testing Strategy

### Backend Tests

- **Unit tests**: Test business logic in isolation
- **Integration tests**: Test DynamoDB operations with LocalStack
- **Mock external dependencies**: AWS services, Cognito, etc.
- **Coverage goal**: Not enforced, but focus on critical paths

### Frontend Tests

- **Component tests**: Test UI behavior with Vitest
- **E2E tests**: Critical user flows with Playwright
- **Mock API calls**: Use MSW (Mock Service Worker)

## Pre-commit Hooks

Pre-commit hooks run automatically before each commit to ensure code quality. When working with Claude:

- **Run pre-commit on changed files**: After Claude modifies files, pre-commit hooks will run automatically on commit
- **Fix issues before committing**: Address any failures from pre-commit checks
- **Manual run**: Use `pre-commit run --files <file1> <file2>` to check specific files
- **Skip if needed**: Use `git commit --no-verify` only when absolutely necessary

Hooks validate:

- Go formatting and linting
- Frontend linting and tests
- Terraform formatting
- Secrets detection
- YAML/JSON syntax

## Code Organization

### Repository Structure

```text
gestione-caselo/
  backend/
    cmd/              # Lambda handlers
    internal/         # Business logic
      graphql/        # GraphQL resolvers
      repository/     # DynamoDB access
      domain/         # Domain models
    go.mod
  frontend/
    src/
      components/
      views/
      router/
      store/
    package.json
  terraform/
    modules/
    environments/
  docs/
    adrs/
  .github/
    workflows/
```

## Development Commands

### Docker Environment

```bash
docker-compose up              # Start all services (LocalStack, backend, frontend)
docker-compose down            # Stop all services
docker-compose logs -f backend # Follow backend logs
docker-compose exec backend sh # Shell into backend container
```

### Backend

```bash
cd backend
go test ./...                    # Run tests
go build -o bootstrap            # Build Lambda binary
GOOS=linux GOARCH=amd64 go build # Build for Lambda deployment

# With LocalStack
AWS_ENDPOINT_URL=http://localhost:4566 go test ./... # Test against LocalStack
```

### Frontend

```bash
cd frontend
npm run dev          # Development server
npm run build        # Production build
npm run test         # Run tests
npm run lint         # Lint code
```

### Infrastructure

```bash
cd terraform/environments/dev
terraform init       # Initialize
terraform plan       # Preview changes
terraform apply      # Apply changes

# With LocalStack (for testing Terraform locally)
tflocal init         # Initialize with LocalStack
tflocal plan         # Plan with LocalStack endpoints
tflocal apply        # Apply to LocalStack
```

### LocalStack

```bash
# Create DynamoDB table locally
aws --endpoint-url=http://localhost:4566 dynamodb create-table \
  --table-name gestione-caselo-dev \
  --attribute-definitions AttributeName=PK,AttributeType=S AttributeName=SK,AttributeType=S \
  --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST

# List local S3 buckets
aws --endpoint-url=http://localhost:4566 s3 ls

# View Cognito user pools
aws --endpoint-url=http://localhost:4566 cognito-idp list-user-pools --max-results 10
```

## Documentation Strategy

- **Check `docs/remote/` first** before asking for help with libraries/tools/languages
- **Create missing docs**: If relevant documentation doesn't exist, search the web and create it
- **Keep docs focused**: Common patterns, gotchas, project-specific usage examples
- **Link to official docs**: Don't duplicate everything, capture what's relevant to this project
- **Structure by technology**: `docs/remote/go/`, `docs/remote/vue/`, `docs/remote/terraform/`, etc.
- **Update when outdated**: Mark docs with last-verified date at the top

Note: `docs/remote/` is gitignored - it's your local knowledge base.

## When I Should Ask You Questions

**Ask me for guidance when**:

- Unfamiliar with Go idioms or patterns (after checking `docs/remote/go/`)
- Unsure about DynamoDB single-table design patterns
- Need clarification on GraphQL resolver implementation (after checking `docs/remote/graphql/`)
- Want examples of testing specific scenarios
- Unclear about Terraform resource configuration (after checking `docs/remote/terraform/`)

**Don't ask me**:

- For encouragement or validation
- To confirm that your approach is "good" (just implement and iterate)
- For permission to proceed with implementation

## Test-First Cycle

Typical workflow for a feature:

1. **I write failing tests** defining the expected behavior
2. **You make tests pass** with minimal implementation
3. **You refactor** if needed (tests still passing)
4. **Repeat** for next feature/test case

If you're stuck, ask for:

- Relevant documentation links
- Explanation of a specific Go/Vue/GraphQL concept
- Examples of similar implementations

## Documentation Expectations

- **Code comments**: Only for non-obvious logic
- **Function docs**: For exported functions (godoc format in Go)
- **ADRs**: For significant architectural decisions
- **README**: Keep updated with setup instructions

## Performance Considerations

- **Lambda cold starts**: Keep dependencies minimal
- **DynamoDB**: Use batch operations where possible
- **Frontend bundles**: Code-split by route, lazy load components
- **CloudFront caching**: Configure appropriate cache headers

## Security Practices

- **Secrets**: Never commit to Git, use environment variables
- **Authentication**: Validate JWT tokens in every Lambda
- **Authorization**: Check user groups for admin operations
- **Input validation**: Validate all user inputs
- **DynamoDB**: Use IAM policies for least-privilege access

## Cost Awareness

This project targets near-zero operating costs:

- Stay within AWS free tier limits
- Use on-demand DynamoDB billing
- Minimize Lambda execution time
- Use CloudFront efficiently

## Deployment Strategy

- **CI/CD**: GitHub Actions handles deployment
- **Environments**: Dev and Production
- **Testing**: All tests must pass before deployment
- **Rollback**: Revert Git commit and redeploy

## Common Patterns

### GraphQL Resolver Pattern

```go
func (r *queryResolver) Booking(ctx context.Context, id string) (*model.Booking, error) {
    // 1. Extract user from context (Cognito claims)
    // 2. Validate authorization
    // 3. Call repository layer
    // 4. Return data or error
}
```

### DynamoDB Repository Pattern

```go
type Repository interface {
    GetBooking(ctx context.Context, id string) (*Booking, error)
    CreateBooking(ctx context.Context, booking *Booking) error
    // ... other methods
}
```

### Vue Component Pattern

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useQuery } from '@vue/apollo-composable'

const data = ref(null)
const { result, loading, error } = useQuery(QUERY)
</script>
```

## Questions Welcome

If something in these rules is unclear or doesn't match your workflow, let me know and we'll adjust.
