import { GraphQLClient } from "graphql-request";

export function useGraphQL() {
  const client = new GraphQLClient(import.meta.env.VITE_API_URL);

  return { client };
}
