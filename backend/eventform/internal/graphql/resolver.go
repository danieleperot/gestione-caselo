package graphql

import (
	"context"

	"github.com/daniele/gestione-caselo/backend/eventform/internal/queue"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type QueueClient interface {
	SendEmailMessage(ctx context.Context, msg queue.EmailMessage) error
}

type Resolver struct {
	queueClient QueueClient
}

func NewResolver(queueClient QueueClient) *Resolver {
	return &Resolver{
		queueClient: queueClient,
	}
}
