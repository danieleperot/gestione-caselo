import { GraphQLClient } from "graphql-request";
import { useAuth } from "./useAuth";

export function useGraphQL() {
  const { getToken } = useAuth();

  const client = new GraphQLClient(import.meta.env.VITE_API_URL, {
    headers: () => {
      const token = getToken();
      const headers: Record<string, string> = {};

      if (token) {
        headers.Authorization = `Bearer ${token}`;
      }

      return headers;
    },
  });

  return { client };
}
