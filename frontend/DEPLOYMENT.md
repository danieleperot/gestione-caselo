# Frontend Deployment Guide

## Runtime Configuration

The frontend uses a runtime configuration file (`public/config.js`) that allows the API URL to be set dynamically at deployment time without rebuilding the application.

### How It Works

1. **Build Time**: The `public/config.js` file contains a placeholder:
   ```javascript
   window.APP_CONFIG = {
     apiUrl: 'PLACEHOLDER_API_URL'
   };
   ```

2. **Deployment Time**: The placeholder is replaced with the actual API Gateway URL using a sed script.

3. **Runtime**: The application reads from `window.APP_CONFIG.apiUrl` if available, otherwise falls back to `VITE_API_URL` for local development.

### Local Development

For local development, the application uses the `VITE_API_URL` environment variable from `.env`:

```bash
# docker-compose already sets this
VITE_API_URL=http://localhost:8080/graphql
```

### CI/CD Deployment

After building the frontend, replace the placeholder before uploading to S3:

```bash
# Build the frontend
npm run build

# Replace the API URL placeholder
./scripts/replace-api-url.sh "https://your-api-gateway-url.execute-api.region.amazonaws.com" "frontend/dist"

# Upload to S3
aws s3 sync frontend/dist s3://your-bucket --delete
```

### GitHub Actions Example

```yaml
- name: Build Frontend
  run: |
    cd frontend
    npm ci
    npm run build

- name: Get API Gateway URL
  id: api-url
  run: |
    API_URL=$(terraform output -raw api_gateway_url)
    echo "url=$API_URL" >> $GITHUB_OUTPUT

- name: Replace API URL
  run: |
    ./scripts/replace-api-url.sh "${{ steps.api-url.outputs.url }}" "frontend/dist"

- name: Deploy to S3
  run: |
    aws s3 sync frontend/dist s3://your-frontend-bucket --delete
```

### Terraform Integration

You can also have Terraform handle this automatically. The API Gateway URL can be obtained from Terraform outputs:

```bash
terraform output -raw api_gateway_url
```

## Verification

After deployment, you can verify the config is correct by visiting:
```
https://yourdomain.com/config.js
```

You should see:
```javascript
window.APP_CONFIG = {
  apiUrl: 'https://actual-api-url.execute-api.region.amazonaws.com'
};
```
