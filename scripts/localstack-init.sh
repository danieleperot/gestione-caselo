#!/bin/bash

echo "Verifying SES domain for local development..."

awslocal ses verify-domain-identity --domain local.test

echo "SES verification complete"
