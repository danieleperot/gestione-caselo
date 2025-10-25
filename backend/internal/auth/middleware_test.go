package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daniele/gestione-caselo/internal/auth"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		wantStatusCode int
		wantUserInCtx  bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			wantStatusCode: http.StatusUnauthorized,
			wantUserInCtx:  false,
		},
		{
			name:           "malformed authorization header",
			authHeader:     "InvalidHeader",
			wantStatusCode: http.StatusUnauthorized,
			wantUserInCtx:  false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.jwt.token",
			wantStatusCode: http.StatusUnauthorized,
			wantUserInCtx:  false,
		},
		// Note: Testing with valid token requires actual JWT from cognito-local
		// This would be an integration test rather than unit test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that checks for user in context
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := auth.GetUserFromContext(r.Context())
				if tt.wantUserInCtx && user == nil {
					t.Error("Expected user in context, got nil")
				}
				if !tt.wantUserInCtx && user != nil {
					t.Errorf("Expected no user in context, got %v", user)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with auth middleware
			config := auth.Config{
				JWKSUrl: "http://localhost:9229/test-pool/.well-known/jwks.json",
			}
			handler := auth.Middleware(config)(testHandler)

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Record response
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.wantStatusCode {
				t.Errorf("Expected status %d, got %d", tt.wantStatusCode, rr.Code)
			}
		})
	}
}
