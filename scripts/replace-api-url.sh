#!/bin/bash

# Script to replace the API URL placeholder in the built frontend
# Usage: ./scripts/replace-api-url.sh <api-gateway-url> <dist-directory>
#
# Example:
#   ./scripts/replace-api-url.sh "https://abc123.execute-api.eu-south-1.amazonaws.com" "frontend/dist"

set -e

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <api-gateway-url> <dist-directory>"
    echo "Example: $0 https://abc123.execute-api.eu-south-1.amazonaws.com frontend/dist"
    exit 1
fi

API_URL="$1"
DIST_DIR="$2"
CONFIG_FILE="$DIST_DIR/config.js"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: config.js not found at $CONFIG_FILE"
    exit 1
fi

echo "Replacing PLACEHOLDER_API_URL with $API_URL in $CONFIG_FILE"
sed -i "s|PLACEHOLDER_API_URL|$API_URL|g" "$CONFIG_FILE"

echo "✓ API URL replacement complete"
