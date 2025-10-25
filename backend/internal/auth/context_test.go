package auth_test

import (
	"context"
	"testing"

	"github.com/daniele/gestione-caselo/internal/auth"
)

func TestUserContext(t *testing.T) {
	t.Run("GetUserFromContext returns nil when no user", func(t *testing.T) {
		ctx := context.Background()
		user := auth.GetUserFromContext(ctx)

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})

	t.Run("GetUserFromContext returns user when present", func(t *testing.T) {
		expectedUser := &auth.User{
			ID:    "test-user-id",
			Email: "test@example.com",
		}

		ctx := auth.SetUserInContext(context.Background(), expectedUser)
		user := auth.GetUserFromContext(ctx)

		if user == nil {
			t.Fatal("Expected user in context, got nil")
		}

		if user.ID != expectedUser.ID {
			t.Errorf("Expected user ID %s, got %s", expectedUser.ID, user.ID)
		}

		if user.Email != expectedUser.Email {
			t.Errorf("Expected user email %s, got %s", expectedUser.Email, user.Email)
		}
	})
}
