#!/usr/bin/env bash

# Script to generate a JWT token from cognito-local for testing
# Usage: ./scripts/get-jwt-token.sh [username] [password]
# Output: JWT token to stdout (all other messages to stderr)

set -e

# Default credentials (test user)
USERNAME="${1:-test@example.com}"
PASSWORD="${2:-Password123!}"

# Load configuration from frontend/.env
ENV_FILE="frontend/.env"

if [ ! -f "$ENV_FILE" ]; then
    echo "Error: frontend/.env file not found" >&2
    echo "Run scripts/setup-cognito-local.sh first" >&2
    exit 1
fi

source "$ENV_FILE"

# Authenticate and get tokens
RESPONSE=$(aws cognito-idp initiate-auth \
    --endpoint-url "$VITE_COGNITO_ENDPOINT" \
    --region "$VITE_COGNITO_REGION" \
    --auth-flow USER_PASSWORD_AUTH \
    --client-id "$VITE_COGNITO_CLIENT_ID" \
    --auth-parameters "USERNAME=$USERNAME,PASSWORD=$PASSWORD" \
    --output json 2>&1)

# Check if authentication was successful
if [ $? -ne 0 ]; then
    echo "Authentication failed!" >&2
    echo "$RESPONSE" >&2
    exit 1
fi

# Extract IdToken
ID_TOKEN=$(echo "$RESPONSE" | jq -r '.AuthenticationResult.IdToken')

if [ "$ID_TOKEN" = "null" ] || [ -z "$ID_TOKEN" ]; then
    echo "Failed to extract IdToken from response" >&2
    echo "$RESPONSE" >&2
    exit 1
fi

# Decode JWT payload and extract expiration time
PAYLOAD=$(echo "$ID_TOKEN" | cut -d '.' -f2)
# Add padding if needed for base64 decoding
PAYLOAD_PADDED=$(echo "$PAYLOAD" | awk '{print $0 (length($0) % 4 == 2 ? "==" : length($0) % 4 == 3 ? "=" : "")}')
EXP_TIMESTAMP=$(echo "$PAYLOAD_PADDED" | base64 -d 2>/dev/null | jq -r '.exp')

if [ "$EXP_TIMESTAMP" != "null" ] && [ -n "$EXP_TIMESTAMP" ]; then
    EXP_DATE=$(date -d "@$EXP_TIMESTAMP" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || date -r "$EXP_TIMESTAMP" '+%Y-%m-%d %H:%M:%S' 2>/dev/null)
    echo "Token expires at: $EXP_DATE" >&2
fi

# Output only the token to stdout
echo "$ID_TOKEN"
