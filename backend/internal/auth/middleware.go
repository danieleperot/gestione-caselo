package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	// Cognito publishes public keys for validation at `/.well-known/.jwks.json`
	JWKSUrl string // URL to fetch public keys for JWT validation
}

type CognitoClaims struct {
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func extractUserFromJWT(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CognitoClaims{}, func(t *jwt.Token) (any, error) {
		// TODO: implement checking against Cognito
		return []byte(""), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CognitoClaims)
	if !ok {
		return nil, errors.New("could not parse JWT")
	}

	return &User{ID: claims.Subject, Email: claims.Email}, nil
}

func Middleware(config Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			user, err := extractUserFromJWT(parts[1])
			if err != nil {
				http.Error(w, "unauthorized token", http.StatusUnauthorized)
				return
			}

			ctx := SetUserInContext(r.Context(), user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
