#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Setting up Cognito Local...${NC}"

# Cognito endpoint
ENDPOINT="http://localhost:9229"

# Create User Pool
echo -e "${BLUE}Creating User Pool...${NC}"
USER_POOL_RESPONSE=$(aws cognito-idp create-user-pool \
  --endpoint-url $ENDPOINT \
  --pool-name gestione-caselo-local \
  --region eu-south-1 \
  --no-cli-pager \
  --output json)

USER_POOL_ID=$(echo $USER_POOL_RESPONSE | jq -r '.UserPool.Id')
echo -e "${GREEN}✓ User Pool created: ${USER_POOL_ID}${NC}"

# Create User Pool Client
echo -e "${BLUE}Creating User Pool Client...${NC}"
CLIENT_RESPONSE=$(aws cognito-idp create-user-pool-client \
  --endpoint-url $ENDPOINT \
  --user-pool-id $USER_POOL_ID \
  --client-name gestione-caselo-client \
  --region eu-south-1 \
  --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH \
  --no-cli-pager \
  --output json)

CLIENT_ID=$(echo $CLIENT_RESPONSE | jq -r '.UserPoolClient.ClientId')
echo -e "${GREEN}✓ Client created: ${CLIENT_ID}${NC}"

# Create test user
echo -e "${BLUE}Creating test user...${NC}"
aws cognito-idp admin-create-user \
  --endpoint-url $ENDPOINT \
  --user-pool-id $USER_POOL_ID \
  --username test@example.com \
  --region eu-south-1 \
  --user-attributes Name=email,Value=test@example.com Name=email_verified,Value=true \
  --message-action SUPPRESS \
  --no-cli-pager \
  --output json > /dev/null

echo -e "${GREEN}✓ Test user created: test@example.com${NC}"

# Set permanent password
echo -e "${BLUE}Setting user password...${NC}"
aws cognito-idp admin-set-user-password \
  --endpoint-url $ENDPOINT \
  --user-pool-id $USER_POOL_ID \
  --username test@example.com \
  --password "Password123!" \
  --region eu-south-1 \
  --permanent \
  --no-cli-pager \
  --output json > /dev/null

echo -e "${GREEN}✓ Password set: Password123!${NC}"

# Save configuration to .env files
echo -e "${BLUE}Saving configuration...${NC}"

# Backend .env
cat > backend/.env <<EOF
COGNITO_ENDPOINT=http://cognito-local:9229
COGNITO_USER_POOL_ID=${USER_POOL_ID}
COGNITO_REGION=eu-south-1
EOF

# Frontend .env
cat > frontend/.env <<EOF
VITE_COGNITO_ENDPOINT=http://localhost:9229
VITE_COGNITO_USER_POOL_ID=${USER_POOL_ID}
VITE_COGNITO_CLIENT_ID=${CLIENT_ID}
VITE_COGNITO_REGION=eu-south-1
EOF

echo -e "${GREEN}✓ Configuration saved to .env files${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Cognito Local Setup Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "User Pool ID: ${BLUE}${USER_POOL_ID}${NC}"
echo -e "Client ID: ${BLUE}${CLIENT_ID}${NC}"
echo -e "Test User: ${BLUE}test@example.com${NC}"
echo -e "Password: ${BLUE}Password123!${NC}"
echo ""
echo -e "Configuration saved to:"
echo -e "  - backend/.env"
echo -e "  - frontend/.env"
