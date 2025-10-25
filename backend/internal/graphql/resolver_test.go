package graphql_test

import (
	"context"
	"testing"

	"github.com/daniele/gestione-caselo/internal/auth"
	"github.com/daniele/gestione-caselo/internal/graphql"
)

func TestHelloQuery(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		wantMessage string
		wantErr     bool
	}{
		{
			name:    "unauthenticated request returns error",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name: "authenticated request returns personalized greeting",
			ctx: auth.SetUserInContext(context.Background(), &auth.User{
				ID:    "user123",
				Email: "test@example.com",
			}),
			wantMessage: "Hello World! How are you doing, test@example.com?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &graphql.Resolver{}

			result, err := resolver.Query().Hello(tt.ctx)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Hello() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Hello() returned error: %v", err)
			}

			if result.Message != tt.wantMessage {
				t.Errorf("Hello() message = %v, want %v", result.Message, tt.wantMessage)
			}
		})
	}
}
