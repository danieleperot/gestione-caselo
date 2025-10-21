# ADR-003: Use GraphQL API

## Status

Accepted

## Date

2025-10-21

## Context

The frontend application needs to communicate with the backend Lambda functions to:

- Query available booking slots by date range
- Create booking requests
- Approve/reject bookings (admin only)
- Manage user profiles
- Query booking history
- Future: Handle payments, notifications, and smart lock codes

We need to design an API layer that:

- Provides efficient data fetching for the booking calendar view
- Supports role-based access control (users vs. administrators)
- Is easy to evolve as new features are added
- Works well with AWS Lambda and API Gateway
- Provides good developer experience for frontend developers

The team has experience with REST APIs but has expressed interest in learning GraphQL.

## Decision

We will implement a **GraphQL API** exposed through AWS API Gateway that routes to Lambda functions.

The API will:

- Use a single GraphQL endpoint for all operations
- Define schema with types for Bookings, Users, Events, TimeSlots
- Implement queries for reading data and mutations for modifications
- Handle authentication via Cognito tokens in request headers
- Use DataLoader pattern for efficient DynamoDB batching

## Consequences

### Positive

- **Flexible Data Fetching**: Frontend requests exactly the data it needs
  - Calendar view can request just date + availability status
  - Detail view can request full booking information with user details
  - Reduces over-fetching and under-fetching
- **Single Endpoint**: All API requests go through one URL
  - Simpler API Gateway configuration
  - Easier to manage CORS and authentication
  - Clear contract between frontend and backend
- **Strongly Typed Schema**: GraphQL schema serves as API documentation
  - Auto-generated TypeScript types for frontend
  - Compile-time validation of queries
  - Excellent IDE autocomplete support
- **Evolutionary Design**: Easy to add new fields without breaking existing clients
  - Adding new optional fields doesn't require API versioning
  - Deprecated fields can be marked and removed gradually
- **Great Developer Experience**:
  - GraphQL Playground for testing queries
  - Introspection enables tooling
  - Clear error messages
- **Efficient for Complex Queries**: Natural fit for booking domain
  - "Get all bookings for October with user details and event information"
  - Single query replaces multiple REST endpoints
- **Learning Opportunity**: Team wants to gain GraphQL experience
  - Increasingly common in modern web development
  - Valuable skill for future projects

### Negative

- **Learning Curve**: Team needs to learn GraphQL concepts
  - Schema definition language
  - Resolvers and data loaders
  - Query optimization and N+1 prevention
  - Mitigation: Good documentation and libraries like gqlgen for Go
- **Complexity for Simple Operations**: Simple CRUD operations require more boilerplate
  - Schema definition + resolver implementation
  - More complex than simple REST endpoint
  - Mitigation: Code generation tools reduce boilerplate
- **Over-fetching Risk**: Poorly designed queries can request too much data
  - Frontend developers might request unnecessary fields
  - Mitigation: Education and query complexity limits
- **Caching Complexity**: HTTP caching is harder than REST
  - Cannot use URL-based CDN caching
  - Mitigation: Not critical for low-traffic booking system; implement application-level caching if needed
- **File Uploads**: GraphQL file upload is less standardized
  - Mitigation: Can use separate REST endpoint or multipart spec
- **Monitoring**: Harder to monitor than REST endpoints
  - Single endpoint makes URL-based metrics useless
  - Mitigation: Log query names and types in structured logs

### Implementation Approach

- Use **gqlgen** for Go (code generation from schema)
- Define schema-first (write GraphQL schema, generate Go code)
- Implement custom resolvers for business logic
- Use DataLoader pattern to batch DynamoDB queries
- Set query complexity limits to prevent abuse
- Enable introspection in development, consider disabling in production

### Example Queries

```graphql
# User booking a slot
mutation CreateBooking {
  createBooking(input: {
    date: "2025-11-15"
    duration: FULL_DAY
    eventDescription: "Birthday party"
    attendees: 50
  }) {
    id
    status
    confirmationCode
  }
}

# Admin viewing bookings
query GetBookings {
  bookings(month: "2025-11", status: PENDING) {
    id
    date
    duration
    user {
      name
      email
    }
    eventDescription
    requestedAt
  }
}

# Calendar availability
query GetAvailability {
  availability(start: "2025-11-01", end: "2025-11-30") {
    date
    morningAvailable
    afternoonAvailable
  }
}
```

## Alternatives Considered

### REST API

- **Pros**:
  - Team has extensive experience
  - Simpler mental model
  - Better HTTP caching
  - Standard CRUD patterns
  - Easier monitoring with URL-based metrics
- **Cons**:
  - Multiple endpoints needed (GET /bookings, POST /bookings, etc.)
  - Over-fetching or under-fetching data
  - Requires API versioning for changes
  - More complex client code with multiple fetch calls
  - Less flexible for complex queries
- **Rejected**: Team wants to learn GraphQL; flexibility benefits outweigh learning curve

### REST API with BFF (Backend for Frontend)

- **Pros**:
  - Optimized endpoints for specific UI views
  - Can evolve mobile/web APIs independently
- **Cons**:
  - More backend code to maintain
  - Duplication across different BFF layers
  - Overkill for single web application
- **Rejected**: Unnecessary complexity for simple booking system

### tRPC (TypeScript RPC)

- **Pros**:
  - End-to-end type safety
  - No schema definition needed
  - Great developer experience for TypeScript projects
- **Cons**:
  - Requires TypeScript on backend (we chose Go)
  - Less mature than GraphQL
  - Smaller ecosystem
- **Rejected**: Not compatible with Go backend choice

### gRPC

- **Pros**:
  - High performance binary protocol
  - Strong typing with Protocol Buffers
  - Bi-directional streaming
- **Cons**:
  - Poor browser support (requires grpc-web proxy)
  - Overkill for low-traffic booking system
  - More complex tooling
  - Less human-readable
- **Rejected**: Web browser requirements make this impractical

## Migration Path

If GraphQL proves too complex or doesn't meet needs:

1. Can add REST endpoints alongside GraphQL
2. Gradually migrate frontend to REST
3. Business logic in resolvers can be reused in REST handlers
4. API Gateway supports both patterns simultaneously

## References

- [GraphQL Official Documentation](https://graphql.org/)
- [gqlgen - Go GraphQL Library](https://gqlgen.com/)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [DataLoader Pattern](https://github.com/graphql/dataloader)
