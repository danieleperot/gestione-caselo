# ADR-004: Use DynamoDB Single-Table Design

## Status

Accepted

## Date

2025-10-21

## Context

The application needs to store several types of entities:

- **Users**: Profile information, contact details, role (user/admin)
- **Bookings**: Booking requests with status (pending, approved, rejected)
- **Events**: Scheduled events in the calendar (both approved bookings and admin-created events)
- **Time Slots**: Availability information for half-day/full-day slots
- Future entities: Payments, Notifications, Audit Logs, Smart Lock Codes

DynamoDB offers two main design approaches:

1. **Multi-Table Design**: Separate table for each entity type (Users table, Bookings table, Events table)
2. **Single-Table Design**: All entities in one table using composite partition/sort keys

Key considerations:

- Cost optimization (DynamoDB pricing based on throughput and storage)
- Query patterns for booking system
- Future scalability requirements
- Development complexity

The application has very low traffic and simple access patterns, but we want to follow DynamoDB best practices.

## Decision

We will use **single-table design** with one DynamoDB table storing all entity types.

The table will use:

- **Partition Key (PK)**: Entity identifier with prefixes (USER#, BOOKING#, EVENT#)
- **Sort Key (SK)**: Related entity or timestamp for sorting
- **GSI1 (Global Secondary Index)**: For querying by different access patterns
- **GSI2**: Reserved for future query patterns

Example key structure:

```text
# User entity
PK: USER#user123
SK: METADATA

# Booking entity
PK: BOOKING#booking456
SK: METADATA

# Event by date (for calendar queries)
PK: EVENT#2025-11-15
SK: MORNING or AFTERNOON

# User's bookings (inverted relationship)
PK: USER#user123
SK: BOOKING#booking456
```

## Consequences

### Positive

- **Cost Efficiency**: Single table with on-demand pricing
  - Free tier: 25GB storage, 2.5M read requests/month
  - Pay-per-request: $0.25 per million reads, $1.25 per million writes
  - Expected cost: $0/month (well within free tier)
  - Multi-table approach would use same resources but with more complexity
- **Better Performance**: Related data can be fetched in single query
  - Get user and all their bookings in one query
  - Get event with booking details in one query
  - Reduces API round trips
- **Simplified Backups**: One table to backup and restore
- **Atomic Transactions**: Can update multiple entity types in single transaction
  - Example: Approve booking + create event + update user stats atomically
- **Follows AWS Best Practices**: Single-table design is recommended pattern by AWS
- **Future-Proof**: Easy to add new entity types without creating new tables
  - Add PAYMENT#, NOTIFICATION#, LOCKCODE# prefixes as needed

### Negative

- **Increased Complexity**: Harder to understand than multi-table design
  - Composite keys require careful planning
  - Schema is less obvious from table structure
  - Mitigation: Document access patterns and key designs clearly
- **Query Limitations**: Cannot easily query across entity types
  - Example: "Get all bookings and all users" requires multiple queries
  - Mitigation: Application doesn't need such queries
- **Development Learning Curve**: Team must learn single-table patterns
  - Understanding partition key design
  - GSI usage and projection
  - Mitigation: Clear examples and documentation
- **Schema Changes More Risky**: Poorly designed keys require migration
  - Changing PK/SK structure requires table rebuild
  - Mitigation: Careful upfront planning of access patterns
- **Local Development**: More complex to visualize and debug
  - Tools like NoSQL Workbench help
  - Mitigation: Good tooling and clear naming conventions

### Access Patterns

Documented access patterns for the booking system:

1. **Get user by ID**: PK = USER#{userId}, SK = METADATA
2. **Get booking by ID**: PK = BOOKING#{bookingId}, SK = METADATA
3. **Get all bookings for user**: PK = USER#{userId}, SK begins_with BOOKING#
4. **Get events by date**: PK = EVENT#{date}, SK = MORNING/AFTERNOON
5. **Get all pending bookings** (admin): GSI1PK = STATUS#PENDING, GSI1SK = {timestamp}
6. **Get availability for date range**: Query multiple EVENT#{date} partitions
7. **Get user by email** (login): GSI1PK = EMAIL#{email}

Future patterns:
8. **Get payments for booking**: PK = BOOKING#{bookingId}, SK begins_with PAYMENT#
9. **Get notification history**: PK = USER#{userId}, SK begins_with NOTIFICATION#
10. **Get lock codes for event**: PK = EVENT#{eventId}, SK = LOCKCODE

### Table Design

```text
Primary Table:
- PK (String): Partition key
- SK (String): Sort key
- GSI1PK (String): Global secondary index 1 partition key
- GSI1SK (String): Global secondary index 1 sort key
- Attributes: Various entity-specific attributes (JSON-like)

Indexes:
- Primary: PK + SK
- GSI1: GSI1PK + GSI1SK (for status queries, email lookups)
- GSI2: Reserved for future use

Capacity Mode: On-demand (pay-per-request)
```

### Entity Examples

```json
// User
{
  "PK": "USER#abc123",
  "SK": "METADATA",
  "GSI1PK": "EMAIL#user@example.com",
  "GSI1SK": "USER",
  "entityType": "User",
  "userId": "abc123",
  "email": "user@example.com",
  "name": "John Doe",
  "phone": "+39 123 456 7890",
  "role": "user",
  "createdAt": "2025-10-21T10:00:00Z"
}

// Booking
{
  "PK": "BOOKING#xyz789",
  "SK": "METADATA",
  "GSI1PK": "STATUS#PENDING",
  "GSI1SK": "2025-10-21T10:00:00Z",
  "entityType": "Booking",
  "bookingId": "xyz789",
  "userId": "abc123",
  "date": "2025-11-15",
  "duration": "FULL_DAY",
  "eventDescription": "Birthday party",
  "attendees": 50,
  "status": "pending",
  "requestedAt": "2025-10-21T10:00:00Z"
}

// User-Booking relationship (for easy querying)
{
  "PK": "USER#abc123",
  "SK": "BOOKING#xyz789",
  "entityType": "UserBooking",
  "bookingId": "xyz789",
  "date": "2025-11-15",
  "status": "pending"
}

// Event
{
  "PK": "EVENT#2025-11-15",
  "SK": "FULLDAY",
  "entityType": "Event",
  "eventId": "evt456",
  "bookingId": "xyz789",
  "date": "2025-11-15",
  "timeSlot": "FULL_DAY",
  "eventType": "private_booking",
  "title": "Birthday party",
  "createdBy": "abc123"
}
```

## Alternatives Considered

### Multi-Table Design (Separate Tables)

- **Pros**:
  - Simpler mental model
  - Clear separation of concerns
  - Easier to understand for new developers
  - Standard relational database thinking
  - Better tooling support
- **Cons**:
  - Same costs as single-table (no savings)
  - More tables to manage and backup
  - Cannot use transactions across tables easily
  - More complex join operations in application code
  - Not following DynamoDB best practices
- **Rejected**: No cost benefit; adds operational complexity

### Relational Database (RDS Aurora Serverless)

- **Pros**:
  - Familiar SQL queries
  - Complex joins and relationships
  - ACID transactions
  - Better tooling and admin interfaces
- **Cons**:
  - Minimum cost ~$0.50/day ($15/month) even with serverless
  - Overkill for simple booking system
  - More complex to manage
  - Slower cold starts
- **Rejected**: 15x higher monthly costs; unnecessary complexity

### Amazon Timestream (Time-Series Database)

- **Pros**:
  - Optimized for time-series data
  - Good for booking calendars
- **Cons**:
  - Designed for high-frequency sensor data
  - Overkill for low-frequency bookings
  - More expensive than DynamoDB
  - Less flexible for user/auth data
- **Rejected**: Wrong use case; optimized for different access patterns

### Amazon DocumentDB (MongoDB-compatible)

- **Pros**:
  - Familiar MongoDB query syntax
  - Flexible document structure
- **Cons**:
  - Not serverless (minimum instance costs)
  - Higher costs than DynamoDB
  - More complex to manage
- **Rejected**: Cost prohibitive for low traffic

## Implementation Guidelines

1. **Document all access patterns upfront** before designing keys
2. **Use descriptive prefixes**: USER#, BOOKING#, EVENT# for clarity
3. **Version entities**: Include version field for optimistic locking
4. **Use ISO 8601 timestamps** for sort keys that need ordering
5. **Implement repository pattern** in Go to abstract DynamoDB details
6. **Use NoSQL Workbench** for testing and visualizing design
7. **Monitor access patterns** and adjust GSI as usage evolves

## References

- [AWS DynamoDB Best Practices](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html)
- [The DynamoDB Book by Alex DeBrie](https://www.dynamodbbook.com/)
- [Single-Table Design with DynamoDB](https://aws.amazon.com/blogs/compute/creating-a-single-table-design-with-amazon-dynamodb/)
- [NoSQL Workbench for DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/workbench.html)
