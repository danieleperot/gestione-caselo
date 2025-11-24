package auth

import "context"

type User struct {
	ID    string
	Email string
}

// Use a private type for context keys to avoid collisions
type contextKey string

const userContextKey contextKey = "user"

func SetUserInContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}

	return user
}
