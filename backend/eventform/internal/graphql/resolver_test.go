package graphql_test

import (
	"context"
	"testing"

	"github.com/daniele/gestione-caselo/backend/eventform/internal/graphql"
	"github.com/daniele/gestione-caselo/backend/eventform/internal/graphql/model"
	"github.com/daniele/gestione-caselo/backend/eventform/internal/queue"
)

type mockQueueClient struct {
	messages []queue.EmailMessage
}

func (m *mockQueueClient) SendEmailMessage(ctx context.Context, msg queue.EmailMessage) error {
	m.messages = append(m.messages, msg)
	return nil
}

func newMockQueueClient() *mockQueueClient {
	return &mockQueueClient{
		messages: make([]queue.EmailMessage, 0),
	}
}

func TestHelloQuery(t *testing.T) {
	resolver := graphql.NewResolver(nil)

	result, err := resolver.Query().Hello(context.Background())

	if err != nil {
		t.Fatalf("Hello() returned error: %v", err)
	}

	expectedMessage := "Hello World!"
	if result.Message != expectedMessage {
		t.Errorf("Hello() message = %v, want %v", result.Message, expectedMessage)
	}
}

func TestSubmitEventBooking(t *testing.T) {
	tests := []struct {
		name    string
		input   model.EventBookingInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid input returns success",
			input: model.EventBookingInput{
				FullName:    "Mario Rossi",
				Email:       "mario.rossi@example.com",
				Phone:       "0123456789",
				Description: "Evento di compleanno",
				Date:        "2025-12-15",
				AcceptData:  true,
				Association: stringPtr("Associazione ricreativa"),
			},
			wantErr: false,
		},
		{
			name: "valid input without optional association",
			input: model.EventBookingInput{
				FullName:    "Mario Rossi",
				Email:       "mario.rossi@example.com",
				Phone:       "0123456789",
				Description: "Evento di compleanno",
				Date:        "2025-12-15",
				AcceptData:  true,
			},
			wantErr: false,
		},
		{
			name: "missing fullName returns error",
			input: model.EventBookingInput{
				Email:       "mario.rossi@example.com",
				Phone:       "0123456789",
				Description: "Evento di compleanno",
				Date:        "2025-12-15",
				AcceptData:  true,
			},
			wantErr: true,
			errMsg:  "fullName is required",
		},
		{
			name: "missing email returns error",
			input: model.EventBookingInput{
				FullName:    "Mario Rossi",
				Phone:       "0123456789",
				Description: "Evento di compleanno",
				Date:        "2025-12-15",
				AcceptData:  true,
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "missing phone returns error",
			input: model.EventBookingInput{
				FullName:    "Mario Rossi",
				Email:       "mario.rossi@example.com",
				Description: "Evento di compleanno",
				Date:        "2025-12-15",
				AcceptData:  true,
			},
			wantErr: true,
			errMsg:  "phone is required",
		},
		{
			name: "missing description returns error",
			input: model.EventBookingInput{
				FullName:   "Mario Rossi",
				Email:      "mario.rossi@example.com",
				Phone:      "0123456789",
				Date:       "2025-12-15",
				AcceptData: true,
			},
			wantErr: true,
			errMsg:  "description is required",
		},
		{
			name: "missing date returns error",
			input: model.EventBookingInput{
				FullName:    "Mario Rossi",
				Email:       "mario.rossi@example.com",
				Phone:       "0123456789",
				Description: "Evento di compleanno",
				AcceptData:  true,
			},
			wantErr: true,
			errMsg:  "date is required",
		},
		{
			name: "acceptData false returns error",
			input: model.EventBookingInput{
				FullName:    "Mario Rossi",
				Email:       "mario.rossi@example.com",
				Phone:       "0123456789",
				Description: "Evento di compleanno",
				Date:        "2025-12-15",
				AcceptData:  false,
			},
			wantErr: true,
			errMsg:  "acceptData must be true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueue := newMockQueueClient()
			resolver := graphql.NewResolver(mockQueue)

			result, err := resolver.Mutation().SubmitEventBooking(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SubmitEventBooking() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("SubmitEventBooking() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("SubmitEventBooking() returned error: %v", err)
			}

			if !result.Success {
				t.Errorf("SubmitEventBooking() success = false, want true")
			}

			if result.Message == "" {
				t.Errorf("SubmitEventBooking() message is empty")
			}

			// Verify that two email messages were queued
			if len(mockQueue.messages) != 2 {
				t.Errorf("Expected 2 email messages to be queued, got %d", len(mockQueue.messages))
			}

			// Verify admin email was queued
			foundAdmin := false
			foundCustomer := false
			for _, msg := range mockQueue.messages {
				if msg.Template == "admin_new_event" {
					foundAdmin = true
				}
				if msg.Template == "customer_new_event" {
					foundCustomer = true
				}
			}

			if !foundAdmin {
				t.Error("admin_new_event email was not queued")
			}
			if !foundCustomer {
				t.Error("customer_new_event email was not queued")
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
