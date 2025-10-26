import { ref } from "vue";
import {
  CognitoUserPool,
  CognitoUser,
  AuthenticationDetails,
  CognitoUserSession,
} from "amazon-cognito-identity-js";
import { usePasswordAuth } from "./local/usePasswordAuth";

export interface User {
  id: string;
  email: string;
}

const userPool = new CognitoUserPool({
  UserPoolId: import.meta.env.VITE_COGNITO_USER_POOL_ID,
  ClientId: import.meta.env.VITE_COGNITO_CLIENT_ID,
  endpoint: import.meta.env.VITE_COGNITO_ENDPOINT,
});

const user = ref<User | null>(null);
const idToken = ref<string | null>(null);

const initalizeAuth = () => {
  const storedUser = localStorage.getItem("user");
  const storedToken = localStorage.getItem("idToken");

  if (storedUser && storedToken) {
    user.value = JSON.parse(storedUser);
    idToken.value = storedToken;
  }
};

initalizeAuth();

export function useAuth() {
  // If VITE_USE_COGNITO_USER_PASSWORD_FLOW is set, forward to password auth
  if (import.meta.env.VITE_USE_COGNITO_USER_PASSWORD_FLOW) {
    return usePasswordAuth();
  }

  const login = (email: string, password: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      const authenticationDetails = new AuthenticationDetails({
        Username: email,
        Password: password,
      });

      const cognitoUser = new CognitoUser({
        Username: email,
        Pool: userPool,
      });

      cognitoUser.authenticateUser(authenticationDetails, {
        onSuccess: (session: CognitoUserSession) => {
          const token = session.getIdToken().getJwtToken();
          const payload = session.getIdToken().payload;

          const userData: User = {
            id: payload.sub,
            email: payload.email,
          };

          localStorage.setItem("user", JSON.stringify(userData));
          localStorage.setItem("idToken", token);

          resolve();
        },
        onFailure: (err) => {
          reject(err);
        },
      });
    });
  };

  const logout = () => {
    user.value = null;
    idToken.value = null;

    localStorage.removeItem("user");
    localStorage.removeItem("idToken");

    const cognitoUser = userPool.getCurrentUser();
    if (cognitoUser) {
      cognitoUser.signOut();
    }
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
