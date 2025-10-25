package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daniele/gestione-caselo/internal/auth"
)

func TestAuthMiddleware(t *testing.T) {
	// Create mock JWKS server
	jwks := newTestJWKS(t)
	defer jwks.close()

	// Generate valid test token
	validToken, err := jwks.createToken("test-user-123", "test@example.com", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

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
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.jwt.token",
			wantStatusCode: http.StatusUnauthorized,
			wantUserInCtx:  false,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			wantStatusCode: http.StatusOK,
			wantUserInCtx:  true,
		},
		{
			name:           "expired token",
			authHeader:     "Bearer " + createExpiredToken(t, jwks),
			wantStatusCode: http.StatusUnauthorized,
			wantUserInCtx:  false,
		},
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

				// If user expected, validate the data
				if tt.wantUserInCtx && user != nil {
					if user.ID != "test-user-123" {
						t.Errorf("Expected user ID 'test-user-123', got '%s'", user.ID)
					}
					if user.Email != "test@example.com" {
						t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
					}
				}

				w.WriteHeader(http.StatusOK)
			})

			// Wrap with auth middleware
			config := auth.Config{
				JWKSUrl: jwks.server.URL + "/.well-known/jwks.json",
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

// Helper to create expired token for testing
func createExpiredToken(t *testing.T, jwks *testJWKS) string {
	token, err := jwks.createToken("expired-user", "expired@example.com", time.Now().Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}
	return token
}
