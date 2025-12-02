#!/bin/bash

echo "Verifying SES domain for local development..."

awslocal ses verify-domain-identity --domain local.test

echo "SES verification complete"

echo "Creating SSM parameters for local development..."

awslocal ssm put-parameter \
  --name "/gestione-caselo/local/ADMIN_EMAILS" \
  --value "admin@local.test,manager@local.test" \
  --type SecureString

awslocal ssm put-parameter \
  --name "/gestione-caselo/local/FROM_ADDRESS" \
  --value "noreply@local.test" \
  --type SecureString

awslocal ssm put-parameter \
  --name "/gestione-caselo/local/COGNITO_ENDPOINT" \
  --value "http://cognito-local:9229" \
  --type SecureString

awslocal ssm put-parameter \
  --name "/gestione-caselo/local/COGNITO_POOL_ID" \
  --value "local_3aR5KXgo" \
  --type SecureString

awslocal ssm put-parameter \
  --name "/gestione-caselo/local/ALLOWED_ORIGINS" \
  --value "http://localhost:5173" \
  --type SecureString

echo "SSM parameters created"
