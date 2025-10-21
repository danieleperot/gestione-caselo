# ADR-005: Use AWS Cognito for Authentication

## Status

Accepted

## Date

2025-10-21

## Context

The application requires user authentication and authorization with two distinct user roles:

- **Regular Users**: Can view calendar, request bookings, see their own bookings
- **Administrators**: Can approve/reject bookings, create events directly, manage all bookings

Authentication requirements:

- User registration with email/password
- Secure login with session management
- Password reset functionality
- Role-based access control (RBAC)
- Token-based API authentication
- Future: Possible MFA for administrators
- Future: Social login (Google, Facebook)

The application is built on AWS serverless architecture and needs a cost-effective solution that integrates well with Lambda and API Gateway.

## Decision

We will use **AWS Cognito User Pools** for authentication and authorization.

Implementation approach:

- Cognito User Pool for user registration and authentication
- Cognito Groups for role management (users, administrators)
- JWT tokens for API authentication
- Custom attributes for storing additional user metadata
- Cognito hosted UI for initial implementation, custom UI later if needed

## Consequences

### Positive

- **Cost Effectiveness**: First 50,000 Monthly Active Users (MAU) are **FREE**
  - Expected usage: 10-100 active users per month
  - **Expected cost: $0/month**
  - Only costs if application exceeds 50,000 MAU (extremely unlikely)
- **Fully Managed Service**: No infrastructure to maintain
  - Handles password hashing, storage, and security
  - Automatic security updates
  - Built-in protection against common attacks (brute force, etc.)
- **Native AWS Integration**: Works seamlessly with other AWS services
  - API Gateway can validate JWT tokens automatically
  - Lambda can access user context from Cognito claims
  - IAM roles can be assumed based on Cognito identity
- **Security Best Practices Built-in**:
  - Bcrypt password hashing
  - Multi-factor authentication support
  - Account recovery workflows
  - Email verification
  - Password complexity requirements
- **OAuth 2.0 and OpenID Connect**: Industry-standard protocols
  - JWT tokens work with any client
  - Easy to integrate with third-party tools
- **User Groups for RBAC**: Built-in group management
  - Create "administrators" and "users" groups
  - Groups included in JWT token claims
  - Simple authorization logic in Lambda functions
- **Rich Feature Set**:
  - Email/SMS verification
  - Password reset via email
  - Social identity providers (Google, Facebook, etc.)
  - Custom attributes (phone, preferences, etc.)
  - Lambda triggers for custom workflows
- **Compliance**: Meets security standards (SOC, PCI DSS, HIPAA eligible)

### Negative

- **Vendor Lock-in**: Tightly coupled to AWS ecosystem
  - Migration to another provider requires significant rewrite
  - User data export/import needed for migration
  - Mitigation: Low traffic means costs remain minimal even long-term
- **Limited Customization**: Some UI/UX limitations
  - Hosted UI is basic (but can use custom UI with SDK)
  - Email templates have constraints
  - Mitigation: Can build custom UI with AWS Amplify or SDKs
- **Learning Curve**: Team needs to learn Cognito concepts
  - User pools, identity pools, app clients
  - JWT token structure and validation
  - Cognito-specific error codes
  - Mitigation: Good documentation available
- **Token Size**: JWT tokens can be large (1-4KB)
  - Includes user attributes and group memberships
  - Not an issue for low-traffic application
  - Mitigation: Only include necessary attributes in tokens
- **Email Sending Limits**: Default limit of 50 emails/day
  - Can be increased or use Amazon SES
  - Mitigation: Low traffic means well within limits
- **Cold Start Overhead**: Token validation adds ~10-50ms to Lambda cold starts
  - Negligible compared to other cold start costs

### Implementation Details

**User Pool Configuration**:

- Sign-in with email and password
- Require email verification
- Password policy: Minimum 8 characters, uppercase, lowercase, numbers
- Allow self-registration
- Account recovery via email

**User Groups**:

```text
Group: administrators
  - Can approve/reject bookings
  - Can create events directly
  - Full access to all API operations

Group: users (default)
  - Can view calendar
  - Can request bookings
  - Can view own bookings
```

**JWT Token Claims**:

```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "cognito:groups": ["users"],
  "cognito:username": "user123",
  "custom:phone": "+39 123 456 7890"
}
```

**Authorization in Lambda**:

```go
func isAdmin(claims map[string]interface{}) bool {
    groups := claims["cognito:groups"].([]string)
    for _, group := range groups {
        if group == "administrators" {
            return true
        }
    }
    return false
}
```

**API Gateway Integration**:

- Use Cognito User Pool Authorizer
- Automatically validates JWT tokens
- Rejects requests with invalid/expired tokens
- Passes user claims to Lambda in event context

### Cost Breakdown

**Cognito User Pools Pricing**:

- First 50,000 MAU: **FREE**
- 50,001 - 100,000 MAU: $0.0055 per MAU
- Advanced security features (optional): $0.05 per MAU

**Additional Costs**:

- SMS MFA: $0.00645 per SMS (only if enabled)
- Email: FREE with Cognito email (50/day limit)
- Email via SES: $0.10 per 1,000 emails (if higher volume needed)

**Expected Monthly Cost**: **$0**

- Estimated active users: 10-100/month
- Well within 50,000 free tier
- No MFA initially
- No SES needed

### Future Enhancements

1. **Multi-Factor Authentication (MFA)**:
   - Enable SMS or TOTP for administrators
   - Cost: ~$0.0065 per SMS

2. **Social Login**:
   - Add Google, Facebook identity providers
   - No additional Cognito cost

3. **Custom Email Templates**:
   - Use Lambda triggers for custom email logic
   - Send via SES for better branding

4. **Advanced Security**:
   - Enable advanced security features for $0.05/MAU
   - Includes compromised credential detection, adaptive authentication

## Alternatives Considered

### Custom Authentication (JWT + DynamoDB)

- **Pros**:
  - Full control over implementation
  - No vendor lock-in
  - Learning experience
- **Cons**:
  - Security risk (rolling your own auth is dangerous)
  - More development time
  - Must implement: password hashing, token generation, email verification, password reset
  - Ongoing maintenance burden
  - Compliance and audit challenges
- **Rejected**: Security concerns; not worth reinventing the wheel

### Auth0

- **Pros**:
  - Excellent developer experience
  - Rich feature set
  - Better UI customization
  - Platform-agnostic
- **Cons**:
  - Costs $23/month for custom domains (Essential plan)
  - Free tier limited to 7,000 active users (vs Cognito's 50,000)
  - Additional service to manage outside AWS
  - Less integration with AWS services
- **Rejected**: Higher costs; prefer AWS-native solution

### Firebase Authentication

- **Pros**:
  - Good developer experience
  - Free tier: 50,000 MAU (same as Cognito)
  - Easy social login
- **Cons**:
  - Google Cloud Platform (different vendor from AWS)
  - Less integration with AWS Lambda/API Gateway
  - Requires cross-cloud setup
  - Team preference for AWS ecosystem
- **Rejected**: Prefer staying within AWS ecosystem

### Amazon Identity Pool (Federated Identities)

- **Pros**:
  - Provides temporary AWS credentials
  - Good for mobile apps accessing AWS directly
- **Cons**:
  - Designed for different use case (direct AWS resource access)
  - Still requires user authentication source (like Cognito User Pool)
  - More complex than needed
- **Rejected**: Not appropriate for web app API authentication

### Supabase Auth

- **Pros**:
  - Open-source
  - Good developer experience
  - Includes database (PostgreSQL)
- **Cons**:
  - Additional service outside AWS
  - Costs $25/month for Pro plan (custom domains, SLAs)
  - Team preference for AWS solutions
- **Rejected**: Prefer AWS-native; higher costs

## Migration Path

If Cognito proves insufficient or too limiting:

1. **Add alternative auth provider** (Auth0, Firebase)
2. **Run both systems in parallel** during migration
3. **Export user data** from Cognito
4. **Import to new provider** with password reset flow
5. **Update Lambda authorizers** to validate new tokens
6. **Deprecate Cognito** once migration complete

Estimated migration time: 1-2 weeks

## Security Considerations

- **Password Policy**: Enforce strong passwords (8+ chars, mixed case, numbers)
- **Email Verification**: Required before account activation
- **Account Lockout**: Built-in protection against brute force
- **Token Expiration**: Access tokens expire after 1 hour (configurable)
- **Refresh Tokens**: Valid for 30 days (configurable)
- **HTTPS Only**: All Cognito endpoints use TLS 1.2+
- **Admin Group Protection**: Manually add users to admin group (no self-service)

## References

- [AWS Cognito Pricing](https://aws.amazon.com/cognito/pricing/)
- [AWS Cognito User Pools Documentation](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools.html)
- [Using Cognito with API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-integrate-with-cognito.html)
- [Cognito JWT Token Validation](https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html)
