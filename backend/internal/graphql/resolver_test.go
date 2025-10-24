package graphql_test

import (
	"context"
	"testing"

	"github.com/daniele/gestione-caselo/internal/graphql"
)

func TestHelloQuery(t *testing.T) {
	resolver := &graphql.Resolver{}

	result, err := resolver.Query().Hello(context.Background())

	if err != nil {
		t.Fatalf("Hello() returned error: %v", err)
	}

	want := "Hello World"
	if result.Message != want {
		t.Errorf("Hello() message = %v, want %v", result.Message, want)
	}
}
