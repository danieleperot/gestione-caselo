import { GraphQLClient } from "graphql-request";

const client = new GraphQLClient(import.meta.env.VITE_API_URL);

export function useGraphQL() {
  return {
    client,
  };
}
