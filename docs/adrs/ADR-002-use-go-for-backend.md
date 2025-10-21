# ADR-002: Use Go for Backend Lambda Functions

## Status

Accepted

## Date

2025-10-21

## Context

AWS Lambda supports multiple programming languages including Python, Node.js, Java, Go, .NET, and Ruby. We need to choose a language for implementing the backend business logic that will handle:

- Booking requests and validation
- Calendar slot availability checks
- User authentication integration
- DynamoDB operations
- Future integrations (payments, notifications, smart locks)

Key considerations:

- Lambda cold start performance (functions idle after periods of inactivity)
- Development speed and maintainability
- Type safety to prevent runtime errors
- AWS SDK support and community resources
- Team learning goals and existing experience

The development team has experience with Python but expressed interest in learning Go for this project.

## Decision

We will use **Go (Golang)** for all backend Lambda functions.

Go is a statically-typed, compiled language developed by Google that is particularly well-suited for serverless architectures.

## Consequences

### Positive

- **Excellent Cold Start Performance**: Go Lambda functions have cold starts of 100-300ms compared to 400-800ms for Python
  - Compiled binaries start faster than interpreted languages
  - Critical for user experience in low-traffic scenarios where cold starts are common
- **Strong Type Safety**: Compile-time type checking catches errors before deployment
  - Reduces runtime errors in production
  - Better IDE support with autocomplete and refactoring
- **Efficient Memory Usage**: Lower memory consumption means lower Lambda costs
  - Can use 128-256MB memory settings for most operations
- **Native Concurrency**: Goroutines enable efficient concurrent operations
  - Useful for future features like batch notifications or parallel API calls
- **Single Binary Deployment**: Compiled executable is easy to deploy
  - No dependency management at runtime
  - Smaller deployment packages than Python with many libraries
- **Learning Opportunity**: Team explicitly wants to learn Go
  - Growing language in cloud-native ecosystems
  - Valuable skill for future projects
- **AWS SDK Support**: Official AWS SDK for Go v2 is well-maintained and documented

### Negative

- **Learning Curve**: Team needs to learn Go syntax, idioms, and best practices
  - Mitigation: Good documentation and simple use cases make learning manageable
  - Booking system complexity is moderate, good for learning
- **Slower Initial Development**: Writing Go code takes more time initially compared to Python
  - Explicit error handling (no exceptions)
  - More verbose than Python
  - Mitigation: Type safety and performance benefits outweigh slower initial development
- **Smaller Ecosystem**: Fewer third-party libraries compared to Python
  - Mitigation: AWS SDK and standard library cover most needs
  - GraphQL libraries (gqlgen) are mature
- **Debugging**: Cannot quickly test snippets in REPL like Python
  - Mitigation: Good testing practices and Go's fast compilation cycle
- **Risk of Switching**: If Go proves too challenging, switching to Python mid-project requires rewriting
  - Mitigation: Team confident in learning Go; Python migration path exists if needed

### Development Workflow

- Use Go modules for dependency management
- Leverage Go's built-in testing framework
- Use AWS SAM or LocalStack for local Lambda testing
- Compilation step required before deployment (handled by CI/CD)

## Alternatives Considered

### Python

- **Pros**:
  - Team has existing experience
  - Faster initial development
  - Large ecosystem of libraries
  - Simpler syntax, less verbose
  - Better for rapid prototyping
- **Cons**:
  - Slower cold starts (400-800ms)
  - Higher memory usage
  - Runtime errors due to dynamic typing
  - Requires dependency packaging
- **Rejected**: Wanted better cold start performance and learning opportunity; team explicitly interested in Go

### Node.js (JavaScript/TypeScript)

- **Pros**:
  - Fast cold starts (200-400ms)
  - Large ecosystem
  - TypeScript adds type safety
  - JSON handling is natural
- **Cons**:
  - Callback hell / async complexity
  - npm dependency management can be complex
  - Team less interested in learning Node.js
- **Rejected**: Go offers better type safety and team preference

### Java/Kotlin

- **Pros**:
  - Strong type safety
  - Mature AWS SDK
  - Spring Framework support
- **Cons**:
  - Very slow cold starts (1-3 seconds)
  - Large deployment packages
  - Heavyweight for simple booking system
- **Rejected**: Cold start performance unacceptable for low-traffic app

### Rust

- **Pros**:
  - Excellent performance and memory safety
  - Fast cold starts
  - Growing AWS Lambda support
- **Cons**:
  - Steep learning curve (ownership, lifetimes)
  - Smaller ecosystem for web development
  - Overkill for booking system complexity
- **Rejected**: Learning curve too steep; Go provides sufficient performance

## Migration Path

If Go proves too challenging or development velocity suffers significantly, we can migrate to Python:

1. Python Lambda functions can coexist with Go functions
2. Migrate one function at a time
3. GraphQL schema remains unchanged (clients unaffected)
4. Estimated migration time: 1-2 weeks for complete rewrite

This provides a safety net while pursuing the Go learning opportunity.

## References

- [AWS Lambda Cold Start Benchmarks](https://mikhail.io/serverless/coldstarts/aws/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [Go Lambda Development](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
