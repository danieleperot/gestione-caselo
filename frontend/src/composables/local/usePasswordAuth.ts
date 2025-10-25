// This file should not be required, HOWEVER cognito-local does not support SRP
// as an authentication flow, even though that's the only flow the JS library
// provided by AWS supports.
//
// Because of this, we are creating this local-only composable that is loaded by
// `useAuth.ts` and replaces its own implementation if the env variable
// VITE_USE_COGNITO_USER_PASSWORD_FLOW is set.
//
// In production we will do things properly and use the library methods instead.
import { ref } from "vue";

export interface User {
  id: string;
  email: string;
}

const COGNITO_ENDPOINT = import.meta.env.VITE_COGNITO_ENDPOINT;
const CLIENT_ID = import.meta.env.VITE_COGNITO_CLIENT_ID;

const user = ref<User | null>(null);
const idToken = ref<string | null>(null);

const initializeAuth = () => {
  const storedUser = localStorage.getItem("user");
  const storedToken = localStorage.getItem("idToken");

  if (storedUser && storedToken) {
    user.value = JSON.parse(storedUser);
    idToken.value = storedToken;
  }
};

initializeAuth();

const decodeJWT = (token: string) => {
  const base64Url = token.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
      .join(""),
  );
  return JSON.parse(jsonPayload);
};

export function usePasswordAuth() {
  const login = async (email: string, password: string): Promise<void> => {
    try {
      const response = await fetch(`${COGNITO_ENDPOINT}/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-amz-json-1.1",
          "X-Amz-Target": "AWSCognitoIdentityProviderService.InitiateAuth",
        },
        body: JSON.stringify({
          ClientId: CLIENT_ID,
          AuthFlow: "USER_PASSWORD_AUTH",
          AuthParameters: {
            USERNAME: email,
            PASSWORD: password,
          },
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Authentication failed");
      }

      const data = await response.json();
      const token = data.AuthenticationResult.IdToken;

      const payload = decodeJWT(token);

      const userData: User = {
        id: payload.sub,
        email: payload.email,
      };

      user.value = userData;
      idToken.value = token;

      localStorage.setItem("user", JSON.stringify(userData));
      localStorage.setItem("idToken", token);
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  };

  const logout = () => {
    user.value = null;
    idToken.value = null;

    localStorage.removeItem("user");
    localStorage.removeItem("idToken");
  };

  const getToken = (): string | null => {
    return idToken.value;
  };

  return {
    user,
    login,
    logout,
    getToken,
  };
}
