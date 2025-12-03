import { GraphQLClient } from "graphql-request";

export function useGraphQL() {
  // Use runtime config if available (for production), otherwise fall back to build-time env var (for local dev)
  const apiUrl = window.APP_CONFIG?.apiUrl !== 'PLACEHOLDER_API_URL'
    ? window.APP_CONFIG?.apiUrl
    : import.meta.env.VITE_API_URL;

  const client = new GraphQLClient(apiUrl);

  return { client };
}
